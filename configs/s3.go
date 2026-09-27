package configs

type S3Configuration struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	Bucket          string
}

func GetS3Configuration() *S3Configuration {
	return &S3Configuration{
		Endpoint:        GetEnv("S3_ENDPOINT", ""),
		AccessKeyID:     GetEnv("S3_ACCESS_KEY_ID", ""),
		SecretAccessKey: GetEnv("S3_SECRET_ACCESS_KEY", ""),
		Region:          GetEnv("S3_REGION", "us-east-1"),
		Bucket:          GetEnv("S3_BUCKET", "noonbyte"),
	}
}
