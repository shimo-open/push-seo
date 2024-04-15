# push-seo 石墨文档开放平台

## 使用

```bash
# 设置环境变量
export AK_ID="YOUR_S3_AK_ID"
export AK_SECRET="YOUR_S3_AK_SECRET"
export BAIDU_TOKEN="YOUR_BAIDU_SEO_TOKEN"

# 执行同步
./push-seo sync -s 'https://host1.com/sitemap.xml' -s 'https://host2.com/sitemap.xml' -d true
```