# concat-gzip-log

concat-gzip-log is a cli tool to concat gzip file on AWS S3 and download.

# background

CloudFront(or API Gateway) generates standard log to S3. Generated log file is compressed by gzip. If you want to analyse log file, you must use analyse tool like Athena. In production environment, you should use Athena to analyse access log. But it is too much cost to use Athena in test environment, staging environment. (If you don't care, you should use Athena. ) So, I made this tool. You can analyse access log file locally, If you this cli.

# how to use

```
A concat cli for gzip access log file on aws s3

Usage:
  concat-gzip-log [flags]

Flags:
  -b, --bucket string   specify target bucket name
  -d, --date string     specify target date YYYY-MM-DD
  -h, --help            help for concat-gzip-log
  -n, --name string     specify output file name (default "access-log.txt")
  -p, --parallel        run parallel
```
