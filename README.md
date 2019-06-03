# Waitter Project (originated in a waiter, responsible for serving)
Waitter (From serving, from waiter) is an agent that performs tasks on files in a configurable way.
With Waitter you can watch directories looking for files (preconfigured) and execute simple tasks like:
- Move, Copy, Delete;
- Zip, UnZip, GZip and GUnzip;
- Change Encoding;
- Upload to AWS S3 or Google Cloud Storage;
- Decrypt .GPG files.
## Configuration
All tasks are defined in a json file, organized by the order of execution.
Example of a simple file move + file delete:
```json
{
  "id": "1",
  "name": "Name of the file.",
  "description": "Description for documentation purposes.",
  "enabled": true,
  "watch": {
    "path": "D:/",
    "regex": "^(example.log)$",
    "recursive": false
  },
  "tasks": [
    {
      "order": 1,
      "type": "FileMove",
      "parameters": [
        {
          "source": "D:/${THIS}$",
          "destination": "D:/back/${THIS}$"
        }
      ]
    },
    {
      "order": 2,
      "type": "FileDelete",
      "parameters": [
        {
          "file": "${THIS}$"
        }
      ]
    }
  ]
}
```
## Graphic Example
### Example 1 - Upload to S3
In the following image we have an example, where several clients / companies upload files to an SFTP server.
On this server, Waitter runs, which detects the new files and performs the processes that have been configured. In one of these cases, Waitter uploads the file to an AWS S3 bucket.
![Example](https://github.com/john-medeiros/waitter/blob/master/doc/doc_example_upload_to_s3.png?raw=true)


