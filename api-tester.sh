#!/bin/bash

if [[ $1 == "http" ]]; then
  PORT=8080
  PROTOCOL=http
else
  PORT=8443
  PROTOCOL=https
fi

API="${PROTOCOL}://localhost:${PORT}/api/v1"

echo "=======Testing API ${API}======="

echo "=======Create 10 users======="
for i in $(seq 1 10);
do
    acct=$(echo $RANDOM | md5sum | head -c 20; echo;)
    fullname="$acct@gmail.com"
    pwd=$(echo $RANDOM | md5sum | head -c 20; echo;)

    result=$(curl -X POST -k -s $API/users \
            -H 'create_user: POST' \
            -H 'Content-Type: application/json' \
            -d '
            {
                "acct": "'$acct'",
                "fullname": "'$fullname'",
                "pwd": "'$pwd'"
            }')
    echo "$result"
done

echo "=======Create last user with same fullname======="
acct=$(echo $RANDOM | md5sum | head -c 20; echo;)
pwd=$(echo $RANDOM | md5sum | head -c 20; echo;)

result=$(curl -X POST -k -s $API/users \
        -H 'create_user: POST' \
        -H 'Content-Type: application/json' \
        -d '
        {
            "acct": "'$acct'",
            "fullname": "'$fullname'",
            "pwd": "'$pwd'"
        }')
echo "$result"

echo "=======User login======="
token=$(curl -X POST -k -s $API/auth/login \
        -H 'user_login: POST' \
        -H 'Content-Type: application/json' \
        -H 'login: POST' \
        -d '
        {
            "acct": "'$acct'",
            "pwd": "'$pwd'"
        }' | jq -r '.token')
echo "token:$token"

echo "=======List users with full name $fullname======="
result=$(curl -X GET -k -s "$API/users?fullname=$fullname" \
        -H 'Content-Type: application/json' \
        -H 'list_user: GET' \
        -H "Authorization: Bearer $token")
echo $result

echo "=======List users with pagination======="
result=$(curl -X GET -k -s "$API/users?limit=3&page=2" \
        -H 'Content-Type: application/json' \
        -H 'list_user: GET' \
        -H "Authorization: Bearer $token")
echo $result

echo "=======List users and get first one for following actions======="
first_id=$(curl -X GET -k -s "$API/users" \
        -H 'Content-Type: application/json' \
        -H 'list_user: GET' \
        -H "Authorization: Bearer $token" | jq -r '.[0].acct')
echo "first user id: $first_id"

echo "=======GET user $first_id"
result=$(curl -X GET -k -s "$API/users/$first_id" \
            -H 'Content-Type: application/json' \
            -H 'get_user: GET' \
            -H "Authorization: Bearer $token")
echo "$result"

echo "=======Update user $first_id ======="
acct=$(echo $RANDOM | md5sum | head -c 20; echo;)
fullname="$acct@gmail.com"
pwd=$(echo $RANDOM | md5sum | head -c 20; echo;)
echo "update fullname:$fullname"
echo "update pwd:$pwd"
result=$(curl -X PUT -k -s "$API/users/$first_id" \
        -H 'Content-Type: application/json' \
        -H 'update_user: PUT' \
        -H "Authorization: Bearer $token" \
        -d '
        {
            "fullname": "'$fullname'",
            "pwd": "'$pwd'"
        }')
echo "$result"

echo "=======Update user full name $first_id ======="
acct=$(echo $RANDOM | md5sum | head -c 20; echo;)
fullname="$acct@gmail.com"
pwd=$(echo $RANDOM | md5sum | head -c 20; echo;)
echo "update fullname:$fullname"
result=$(curl -X PUT -k -s "$API/users/$first_id/fullname" \
        -H 'Content-Type: application/json' \
        -H 'update_user: PUT' \
        -H "Authorization: Bearer $token" \
        -d '
        {
            "fullname": "'$fullname'"
        }')
echo "$result"

echo "=======DELETE user $first_id======="
result=$(curl -X DELETE -k -s "$API/users/$first_id" \
        -H 'Content-Type: application/json' \
        -H 'delete_user: DELETE' \
        -H "Authorization: Bearer $token")
echo "$result"
