package storage

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	"google.golang.org/appengine"
)

type Bucket struct {
	ctx context.Context
}

func NewBucket(r *http.Request) {

}

func (p *Bucket) test(r *http.Request) {
	p.ctx = appengine.NewContext(r)

	projectID := "alert-tempo-159112"
	if projectID == "" {
		fmt.Fprintf(os.Stderr, "GOOGLE_CLOUD_PROJECT environment variable must be set.\n")
		os.Exit(1)
	}

	// [START setup]
	client, err := storage.NewClient(p.ctx)
	if err != nil {
		log.Fatal(err)
	}
	// [END setup]

	// Give the bucket a unique name.
	name := fmt.Sprintf("golang-example-buckets-%d", time.Now().Unix())
	if err := p.create(client, projectID, name); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("created bucket: %v\n", name)

	// list buckets from the project
	buckets, err := p.list(client, projectID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("buckets: %+v\n", buckets)

	// delete the bucket
	if err := p.delete(client, name); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("deleted bucket: %v\n", name)
}

func (p *Bucket) create(client *storage.Client, projectID, bucketName string) error {
	// [START create_bucket]
	if err := client.Bucket(bucketName).Create(p.ctx, projectID, nil); err != nil {
		return err
	}
	// [END create_bucket]
	return nil
}

func (p *Bucket) createWithAttrs(client *storage.Client, projectID, bucketName string) error {
	// [START create_bucket_with_storageclass_and_location]
	bucket := client.Bucket(bucketName)
	if err := bucket.Create(p.ctx, projectID, &storage.BucketAttrs{
		StorageClass: "COLDLINE",
		Location:     "asia",
	}); err != nil {
		return err
	}
	// [END create_bucket_with_storageclass_and_location]
	return nil
}

func (p *Bucket) list(client *storage.Client, projectID string) ([]string, error) {
	// [START list_buckets]
	var buckets []string
	it := client.Buckets(p.ctx, projectID)
	for {
		battrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		buckets = append(buckets, battrs.Name)
	}
	// [END list_buckets]
	return buckets, nil
}

func (p *Bucket) delete(client *storage.Client, bucketName string) error {
	// [START delete_bucket]
	if err := client.Bucket(bucketName).Delete(p.ctx); err != nil {
		return err
	}
	// [END delete_bucket]
	return nil
}

func (p *Bucket) writeBucket(client *storage.Client, bucket, object string) error {
	// [START upload_file]
	f, err := os.Open("notes.txt")
	if err != nil {
		return err
	}
	defer f.Close()

	wc := client.Bucket(bucket).Object(object).NewWriter(p.ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}
	// [END upload_file]
	return nil
}

func (p *Bucket) listBucket(client *storage.Client, bucket string) error {
	ctx := context.Background()
	// [START storage_list_files]
	it := client.Bucket(bucket).Objects(ctx, nil)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Println(attrs.Name)
	}
	// [END storage_list_files]
	return nil
}

func (p *Bucket) listByPrefix(client *storage.Client, bucket, prefix, delim string) error {
	ctx := context.Background()
	// [START storage_list_files_with_prefix]
	// Prefixes and delimiters can be used to emulate directory listings.
	// Prefixes can be used filter objects starting with prefix.
	// The delimiter argument can be used to restrict the results to only the
	// objects in the given "directory". Without the delimeter, the entire  tree
	// under the prefix is returned.
	//
	// For example, given these blobs:
	//   /a/1.txt
	//   /a/b/2.txt
	//
	// If you just specify prefix="/a", you'll get back:
	//   /a/1.txt
	//   /a/b/2.txt
	//
	// However, if you specify prefix="/a"" and delim="/", you'll get back:
	//   /a/1.txt
	it := client.Bucket(bucket).Objects(ctx, &storage.Query{
		Prefix:    prefix,
		Delimiter: delim,
	})
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Println(attrs.Name)
	}
	// [END storage_list_files_with_prefix]
	return nil
}

func (p *Bucket) readBucket(client *storage.Client, bucket, object string) ([]byte, error) {
	ctx := context.Background()
	// [START download_file]
	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	return data, nil
	// [END download_file]
}

func (p *Bucket) attrsBucket(client *storage.Client, bucket, object string) (*storage.ObjectAttrs, error) {
	ctx := context.Background()
	// [START get_metadata]
	o := client.Bucket(bucket).Object(object)
	attrs, err := o.Attrs(ctx)
	if err != nil {
		return nil, err
	}
	return attrs, nil
	// [END get_metadata]
}

func (p *Bucket) makePublicBucket(client *storage.Client, bucket, object string) error {
	ctx := context.Background()
	// [START public]
	acl := client.Bucket(bucket).Object(object).ACL()
	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return err
	}
	// [END public]
	return nil
}

func (p *Bucket) move(client *storage.Client, bucket, object string) error {
	ctx := context.Background()
	// [START move_file]
	dstName := object + "-rename"

	src := client.Bucket(bucket).Object(object)
	dst := client.Bucket(bucket).Object(dstName)

	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
		return err
	}
	if err := src.Delete(ctx); err != nil {
		return err
	}
	// [END move_file]
	return nil
}

func (p *Bucket) copyToBucket(client *storage.Client, dstBucket, srcBucket, srcObject string) error {
	ctx := context.Background()
	// [START copy_file]
	dstObject := srcObject + "-copy"
	src := client.Bucket(srcBucket).Object(srcObject)
	dst := client.Bucket(dstBucket).Object(dstObject)

	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
		return err
	}
	// [END copy_file]
	return nil
}

func (p *Bucket) deleteBucket(client *storage.Client, bucket, object string) error {
	ctx := context.Background()
	// [START delete_file]
	o := client.Bucket(bucket).Object(object)
	if err := o.Delete(ctx); err != nil {
		return err
	}
	// [END delete_file]
	return nil
}

func (p *Bucket) writeEncryptedObject(client *storage.Client, bucket, object string, secretKey []byte) error {
	ctx := context.Background()

	// [START storage_upload_encrypted_file]
	obj := client.Bucket(bucket).Object(object)
	// Encrypt the object's contents.
	wc := obj.Key(secretKey).NewWriter(ctx)
	if _, err := wc.Write([]byte("top secret")); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}
	// [END storage_upload_encrypted_file]
	return nil
}

func (p *Bucket) readEncryptedObject(client *storage.Client, bucket, object string, secretKey []byte) ([]byte, error) {
	ctx := context.Background()

	// [START storage_download_encrypted_file]
	obj := client.Bucket(bucket).Object(object)
	rc, err := obj.Key(secretKey).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	// [END storage_download_encrypted_file]
	return data, nil
}

func (p *Bucket) rotateEncryptionKey(client *storage.Client, bucket, object string, key, newKey []byte) error {
	ctx := context.Background()
	// [START storage_rotate_encryption_key]
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	obj := client.Bucket(bucket).Object(object)
	// obj is encrypted with key, we are encrypting it with the newKey.
	_, err = obj.Key(newKey).CopierFrom(obj.Key(key)).Run(ctx)
	if err != nil {
		return err
	}
	// [END storage_rotate_encryption_key]
	return nil
}
