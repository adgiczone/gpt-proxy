data='{
    "url": "https://api.openai.com/v1/chat/completions",
    "content": "{\"model\": \"gpt-3.5-turbo\",\"messages\": [{\"role\": \"system\",\"content\": \"You are a poetic assistant, skilled in explaining complex programming concepts with creative flair.\"},{\"role\": \"user\",\"content\": \"Compose a poem that explains the concept of recursion in programming.\"}]}"
}'

curl https://www.nekopadia.top/proxy -H "Content-Type: application/json" -H "Authorization: Bearer asdfasdf" -d "$data"

# curl --noproxy '*' http://localhost:8080/proxy -H "Content-Type: application/json" -H "Authorization: Bearer asdfasdf" -d "$data"