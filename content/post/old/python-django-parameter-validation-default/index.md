---
title: "扩展 Django forms.Field - 支持 default 属性"
date: 2022-09-05T22:22:13+08:00
draft: false
reward: false
categories:
 - 后端技术
 - Python
tags:
 - Django
 - 参数验证
 - Python
---

`django.forms` 包提供了 HTML 表单验证的功能，在没有使用 DRF 的情况下，无法合理地处理 API 传参的验证，其中传参验证中就缺少了参数默认值的设置。

## 层次结构

`django.forms` 包提供的 `Field` 类如下：

```python
__all__ = (
    'Field', 'CharField', 'IntegerField',
    'DateField', 'TimeField', 'DateTimeField', 'DurationField',
    'RegexField', 'EmailField', 'FileField', 'ImageField', 'URLField',
    'BooleanField', 'NullBooleanField', 'ChoiceField', 'MultipleChoiceField',
    'ComboField', 'MultiValueField', 'FloatField', 'DecimalField',
    'SplitDateTimeField', 'GenericIPAddressField', 'FilePathField',
    'JSONField', 'SlugField', 'TypedChoiceField', 'TypedMultipleChoiceField',
    'UUIDField',
)
```

继承关系如下图：

```bash
Field: 字段基类
- InlineForeignKeyField: 关联模型的主键字段
- CharField: 字符串字段
	- RegexField: 正则字符串字段
	- UrlField: URL 字段
	- SlugField: 是否允许 unicode 字符串
	- GenericIPAdressField: IP 字段
	- EmailField: 邮箱字段
	- UUIDField: UUID 字段
	- JSONField: JSON 字符串
- IntegerField: 整数字段
	- FloatField: 浮点数字段
	- DecimalField: 十进制小数字段
- BaseTemporalField: 表示时间
	- DatelField: 日期字段
	- TimeField: 时间字段
	- DateTimeField: （日期、时间）字段
- DurationField: 时间段字段
- FileField: 文件字段
- ImageField: 图片字段
- BooleanField: 布尔字段
- NullBooleanField: 布尔字段（可为 null）
- ChoiceField: 可选项字段
	- ModelChoiceField: 模型选项字段（QuerySet）
	- ModelMultipleChoiceField: 模型选项字段（QuerySet，可多个）
	- TypedChoiceField: 可选项字段（强制类型）
	- MultipleChoiceField: 可选项字段（多选）
	- TypedMultipleChoiceField: 可选项字段（多选，强制类型）
	- FilePathField: 文件路径字段
- ComboField: 组合字段，同时使用多个字段验证器
- MultiValueField: 聚合多个字段的逻辑
	- SplitDateTimeField: 多个时间值字段
```

基类 `Field` 的构造函数的声明如下：

```python
def __init__(self, *, required=True, widget=None, label=None, initial=None,
                 help_text='', error_messages=None, show_hidden_initial=False,
                 validators=(), localize=False, disabled=False, label_suffix=None):
```

### initial 参数

其中有一个 `initial` 参数，可以用来表示 `Form` 渲染出的 HTML 代码中的显示值，但并不会给字段赋值，如：

```python
from django.http import HttpRequest, JsonResponse
from django import forms


class ExampleForm(forms.Form):
    name = forms.CharField(required=False, initial='a2htray')


def django_forms_field_initial(request: HttpRequest):
    form = ExampleForm(request.GET)
    if form.is_valid():
        data = form.cleaned_data
        print('data', data)
        print(form)
    else:
        print(form.errors)
        return JsonResponse({})

    return JsonResponse({})
```

请求该接口打印的信息如下：

```bash
data {'name': ''}
<tr><th><label for="id_name">Name:</label></th><td><input type="text" name="name" id="id_name"></td></tr>
```

可见，django 提供的 `Field` 类没有提供类似于默认值的功能。比如，当用户不能给出默认的 `name` 值时进行赋值。

## 自定义 default 参数

封装一个需要提供默认值功能的字段类，这里以 `TypedChoiceField` 为例，其中 `TypedChoiceField._coerce` 方法本身也提供了一定的验证功能，看源码：

```python
class TypedChoiceField(ChoiceField):
	# ...

    def _coerce(self, value):
        if value == self.empty_value or value in self.empty_values:
            return self.empty_value  # 需要的修改地方
        try:
            value = self.coerce(value)
        except (ValueError, TypeError, ValidationError):
            raise ValidationError(
                self.error_messages['invalid_choice'],
                code='invalid_choice',
                params={'value': value},
            )
        return value

    def clean(self, value):
        value = super().clean(value)
        return self._coerce(value)
```

当传值解析为指定的 `empty_value` 或在 `empty_values` 中时，返回指定的默认值，其余不变。

```python
class MyTypeChoiceField(forms.TypedChoiceField):
    def __init__(self, *, default=None, **kwargs):
        self.default = default
        super().__init__(**kwargs)

    def _coerce(self, value):
        if value == self.empty_value or value in self.empty_values:
            return self.default
        else:
            try:
                value = self.coerce(value)
            except (ValueError, TypeError, ValidationError):
                raise ValidationError(
                    self.error_messages['invalid_choice'],
                    code='invalid_choice',
                    params={'value': value},
                )
        return value
```

那么在使用 `MyTypedChoiceField` 时可以指定 `default` 参数，如：

```python
GENDER_CHOICES = (
    (0, 'Male'),
    (1, 'Female'),
)

class Payload(forms.Form):
    gender = MyTypeChoiceField(coerce=lambda p: int(p), choices=GENDER_CHOICES, required=False, default=0)

def get_payload_validate(request: HttpRequest):
    payload = Payload(request.GET)
    if not payload.is_valid():
        return JsonResponse({
            'errors': payload.errors,
        })
    payload = payload.cleaned_data

    return JsonResponse(payload)
```

不带任何参数值请求该接口时的返回如下：

```json
{
    "gender": 0
}
```

不同的 `Field` 子类也可进行类似的继承再修改。