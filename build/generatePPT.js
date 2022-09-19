const fs = require('fs')
const { exec } = require("child_process")

const dir = 'content/posts'

fs.readdir(dir, (error, files) => {
    if (error) {
        console.error(error)
    } else {
        let folders = files.filter(file => file.startsWith('ppt'))
        for (let i = 0; i < folders.length; i++) {
            let folder = folders[0]
            console.log(`process ${folder}`)

            let input = `${dir}/${folder}/index.md`
            let output = `docs/ppt/${folder}.html`
            let cmd = `marp ${input} -o ${output}`
            console.log(`run command: ${cmd}`)

            let child = exec(cmd, function(error, stdout, stderr) {
                console.log('stdout: ' + stdout)
                console.log('stderr: ' + stderr)
                if (error !== null) {
                    console.log('exec error: ' + error)
                }
            })

            child.stdout.pipe(process.stdout)
            child.on('exit', function() {
                process.exit()
            })
        }
    }
})