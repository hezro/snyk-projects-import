# snyk-projects-import

This Go utility allows you import a single file, multiple files, or the whole repo.


***

> Note: This util is a proof of concept for how you can use the [Snyk API](https://snyk.docs.apiary.io/#reference/import-projects) to import files and repos into Snyk. 


# Usage
By default the whole repo will be imported if the `file_path` flag is not specified.
- `--help` - show help and all available flags
- `--token=` - Snyk API token
- `--gitId=` - Git integration ID: This can be found on the Integration page in the Settings area for all integrations that have been configured. https://app.snyk.io/org/*org-name*/manage/integrations Each org will have a unique integration ID for Git
- `--orgId=` - Snyk Target/Destination Organization ID
- `--owner` - For Github: account owner of the repository, For Azure Repos, this is Project ID, For Bitbucket Cloud, this is the Workspace ID
- `--repoName=` - Name of the Repo
- `--branchName=` - Name of the Branch
- `--filePath=` - Optional: Relative path to one or more files. If passing in more than one file use a comma `,` to separate. e.g., `--filePath=Dockerfile,kubernetes/goof-deployment.yaml`
- `@args` - Optional: You can also read one or more arguments from a file. 

- ** Flags or arguments cannot be repeated. If you have a flag defined in the args file then you cannot also pass that same flag in on the command-line. **


# Examples
You can find usage instructions by running:

```bash
snyk-projects-import --help
```

How to import or refresh a full repo:
```bash
snyk-projects-import --token=<Snyk Token> --gitId=<git integration id> --orgId=<organization id> --owner=<repo owner> --repoName=<repo name> --branchName=<branch name>
```


Use the args files to pass in the Snyk Token:
```bash
echo --token=<snyk token> > args
```
```bash
snyk-projects-import @args --gitId=<git integration id> --orgId=<organization id> --owner=<repo owner> --repoName=<repo name> --branchName=<branch name>
```


How to import a Dockerfile in the root of the repo:
```bash
snyk-projects-import --token=<Snyk Token> --gitId=<git integration id> --orgId=<organization id> --owner=<repo owner> --repoName=<repo name> --branchName=<branch name> --filePath=Dockerfile
```

How to import multiple files that exist in a directory:
```bash
snyk-projects-import --token=<Snyk Token> --gitId=<git integration id> --orgId=<organization id> --owner=<repo owner> --repoName=<repo name> --branchName=<branch name> --filePath=kubernetes/goof-deployment.yaml, kubernetes/goof-mongo-deployment.yaml
```

Pass in all flags/arguments via the args file:
```bash
snyk-projects-import @args
```
Args file example:

```bash
--token=11111-1111-1111-1111-111111111111
--gitId=11111-2222-3333-4444-555555555555
--orgId=3333-3333-3333-3333-33333333333
--owner=hezro
--repoName=goof
--branchName=main
```
