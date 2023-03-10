# s3-bucket-intelligence-tiering-compliance
applies intelligent tiering to all s3 buckets in account to save on s3 cost without having to delete objects.

## What is this for ?
Working with legacy s3 buckets with large amounts of data can be tedious, In order to identify resources that are infrequently accessed you can apply intelligent tiering to help locate and save money on stored objects.

## How it works
You can run this one time, but recomended to use as a lambda that runs ones a month in each aws account. It will apply the policy after one day and begin the Intelligent tiering process to save money on s3 storage based on accesability. After a year of no activity to the object it will move it to a long term storage which reduces cost by up to 90%. If needed the object can be restored.

## Benefits

- allows long term solution for a significantly cheaper price than using standard s3
- helps identify and monitor bucket usage
- allows older objects to be retained without having to delete in case object is infrequently accessed compared to normal lifecycles

## Reponse syntax
Once you run it you can check the policy on the bucket and will also return the policy applied in the logs.

![example bucket](assets/policy.png)

```
INFO[0001] Intelligent Tiering applied to s3bucket_name
INFO[0001] {
  Id: "DeepArchive365",
  Status: "Enabled",
  Tierings: [{
      AccessTier: "DEEP_ARCHIVE_ACCESS",
      Days: 365
    }]
}
```

## Docs
[Intelligent Tiering](https://aws.amazon.com/s3/storage-classes/intelligent-tiering/)

### Testing

run `go test -v` in the main directory to generate test results
