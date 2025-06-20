```bash
# steam quotes
watch -n 1 'curl localhost:3001/v1/quotes | jq'

# set next quote
curl -X POST localhost:3001/v1/setNextQuote -d '{"nextQuote": 200}'
```
