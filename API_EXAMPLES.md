# 📡 API Examples - MMS Backend

Exemples d'utilisation de l'API avec curl.

## 🔐 Authentication

### Signup (Inscription)

```bash
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "email": "alice@example.com",
    "password": "Alice1234!",
    "phone": "+261340000001",
    "language": "fr"
  }'
```

**Response:**
```json
{
  "message": "Compte créé avec succès",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": "uuid",
      "username": "alice",
      "email": "alice@example.com",
      "language": "fr"
    }
  }
}
```

### Login (Connexion)

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "alice@example.com",
    "password": "Alice1234!"
  }'
```

**Note:** `identifier` peut être un email, username ou phone.

### Get Current User

```bash
curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 👥 Users

### List Users

```bash
curl "http://localhost:8080/api/v1/users?limit=20&offset=0" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Search Users

```bash
curl "http://localhost:8080/api/v1/users/search?q=alice&limit=10" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Get User by ID

```bash
curl http://localhost:8080/api/v1/users/USER_UUID \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 💬 Direct Messages

### Send Message

```bash
curl -X POST http://localhost:8080/api/v1/messages \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "receiver_id": "uuid-of-receiver",
    "content": "Salut! Comment ça va?"
  }'
```

### Get Conversation

```bash
curl "http://localhost:8080/api/v1/messages/conversation/USER_UUID?limit=50&offset=0" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Get Recent Conversations

```bash
curl "http://localhost:8080/api/v1/messages/conversations?limit=20" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Mark Messages as Read

```bash
curl -X PUT http://localhost:8080/api/v1/messages/read/SENDER_UUID \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Get Unread Count

```bash
curl http://localhost:8080/api/v1/messages/unread/count \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 👥 Groups

### Create Group

```bash
curl -X POST http://localhost:8080/api/v1/groups \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Amis Proches",
    "description": "Notre groupe d'\''amis",
    "type": "private",
    "member_ids": ["uuid1", "uuid2"]
  }'
```

### Get My Groups

```bash
curl http://localhost:8080/api/v1/groups/my \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Get Group Details

```bash
curl http://localhost:8080/api/v1/groups/GROUP_UUID \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Send Group Message

```bash
curl -X POST http://localhost:8080/api/v1/groups/messages \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "group_id": "uuid-of-group",
    "content": "Salut tout le monde!"
  }'
```

### Get Group Messages

```bash
curl "http://localhost:8080/api/v1/groups/GROUP_UUID/messages?limit=50&offset=0" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Delete Group

```bash
curl -X DELETE http://localhost:8080/api/v1/groups/GROUP_UUID \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 🌐 WebSocket

### Connect to WebSocket

**JavaScript Example:**
```javascript
const token = "YOUR_JWT_TOKEN";
const ws = new WebSocket("ws://localhost:8080/api/v1/ws");

// Add Authorization header (if supported by client)
// Or connect with: ws://localhost:8080/api/v1/ws?token=${token}

ws.onopen = () => {
  console.log("Connected to WebSocket");
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log("Received:", message);
};

ws.onerror = (error) => {
  console.error("WebSocket error:", error);
};
```

### Send Direct Message via WebSocket

```javascript
ws.send(JSON.stringify({
  type: "new_message",
  receiver_id: "uuid-of-receiver",
  content: "Hello via WebSocket!"
}));
```

### Send Group Message via WebSocket

```javascript
ws.send(JSON.stringify({
  type: "new_group_message",
  group_id: "uuid-of-group",
  content: "Hello group!"
}));
```

### Send Typing Indicator

```javascript
ws.send(JSON.stringify({
  type: "typing",
  receiver_id: "uuid-of-receiver"
}));
```

### Heartbeat/Ping

```javascript
ws.send(JSON.stringify({
  type: "ping"
}));
```

### Events Received from Server

```javascript
ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  
  switch(data.type) {
    case "user_joined":
      console.log(`${data.data.username} joined`);
      break;
    case "user_left":
      console.log(`${data.data.username} left`);
      break;
    case "new_message":
      console.log("New message:", data);
      break;
    case "new_group_message":
      console.log("New group message:", data);
      break;
    case "pong":
      console.log("Server is alive");
      break;
  }
};
```

---

## 🧪 Testing with Variables

**Set Token Variable (Bash):**
```bash
# Login and save token
TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"identifier":"alice@example.com","password":"Alice1234!"}' \
  | jq -r '.data.token')

# Use the token
curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer $TOKEN"
```

**PowerShell:**
```powershell
# Login and save token
$response = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/login" `
  -Method POST `
  -Body (@{identifier="alice@example.com"; password="Alice1234!"} | ConvertTo-Json) `
  -ContentType "application/json"
  
$token = $response.data.token

# Use the token
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/me" `
  -Headers @{Authorization="Bearer $token"}
```

---

## 📋 Complete Workflow Example

```bash
# 1. Signup Alice
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","email":"alice@test.com","password":"Test1234!"}'

# 2. Signup Bob
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"username":"bob","email":"bob@test.com","password":"Test1234!"}'

# 3. Alice logs in
ALICE_TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"identifier":"alice@test.com","password":"Test1234!"}' \
  | jq -r '.data.token')

# 4. Bob logs in
BOB_TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"identifier":"bob@test.com","password":"Test1234!"}' \
  | jq -r '.data.token')

# 5. Get Bob's ID
BOB_ID=$(curl -s "http://localhost:8080/api/v1/users/search?q=bob" \
  -H "Authorization: Bearer $ALICE_TOKEN" \
  | jq -r '.data[0].id')

# 6. Alice sends message to Bob
curl -X POST http://localhost:8080/api/v1/messages \
  -H "Authorization: Bearer $ALICE_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"receiver_id\":\"$BOB_ID\",\"content\":\"Hello Bob!\"}"

# 7. Bob checks his messages
curl "http://localhost:8080/api/v1/messages/conversation/$ALICE_ID" \
  -H "Authorization: Bearer $BOB_TOKEN"
```

---

## 🔒 Security Notes

- **Toujours utiliser HTTPS** en production
- **Ne jamais exposer** votre JWT token
- **Rotation des tokens** recommandée
- **Rate limiting** recommandé pour la production
- **Validation** de tous les inputs côté serveur

---

## 🐛 Error Responses

```json
// Unauthorized (401)
{
  "error": "invalid or expired token"
}

// Bad Request (400)
{
  "error": "validation error message"
}

// Not Found (404)
{
  "error": "resource not found"
}

// Internal Server Error (500)
{
  "error": "internal server error"
}
```

---

**📖 Plus d'infos**: Voir [README.md](README.md) pour la documentation complète.

