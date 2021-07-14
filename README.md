# cloudfront_remover
Disable and then delete CloudFront Distributions and OAIs

Binaries for Windows, MacOS, and Linux are available on the [Releases](https://github.com/jftuga/cloudfront_remover/releases) page.

## Usage
```
cloudfront_remover
Disable and then delete AWS CloudFront Distributions and associated OAIs

Set the AWS_PROFILE environment variable to use a different profile from the AWS credential file.

Usage:
  cloudfront_remover [command]

Available Commands:
  delete      Delete a CloudFront Distribution
  deleteOAI   Delete CloudFront Origin Access Identifier (OAI)
  disable     Disable a CloudFront Distribution
  list        List distributions and their OAIs
  listOAI     List CloudFront Origin Access Identities (OAIs)
  s3search    Search for an OAI in all S3 bucket policies
```

## Examples

* List distributions and OAIs
```
$ cloudfront list

+----------------+----------------+---------+----------+-------------------+--------------------------------------------------+--------------------------+
|       ID       |      ETAG      | ENABLED |  STATUS  |     1ST ALIAS     |                     1ST OAI                      |         COMMENT          |
+----------------+----------------+---------+----------+-------------------+--------------------------------------------------+--------------------------+
| E012345678912  | E0123456789123 | true    | Deployed | www.example.io    | N/A                                              | N/A                      |
| E987654321098  | E9876543210987 | true    | Deployed | www.example.com   | origin-access-identity/cloudfront/E9876543210987 | terraform--example.com   |
+----------------+----------------+---------+----------+-------------------+--------------------------------------------------+--------------------------+
```

* Delete a distribution
```
PS C:\> .\cloudfront_remover.exe delete -i E012345678912
(no output upon success)
```

* Search all S3 Buckets for an OAI (with region hint)
```
$ cloudfront_remover.exe s3search -i E012345678912 -r us-east-1

+-----------------------+-----------+---------------+
|        BUCKET         |  REGION   |      OAI      |
+-----------------------+-----------+---------------+
| www.stadiumscores.com | us-east-1 | E012345678912 |
+-----------------------+-----------+---------------+
```

## Acknowledgment
* [AWS GoLang SDK](https://aws.amazon.com/sdk-for-go/)
* [ASCII table in golang](github.com/olekukonko/tablewriter)
