# securechat

This repo implements a simple secure chat with two users over tls websocket. The server is [GO code](https://github.com/gorilla/websocket/tree/efaec3cbd167c850a8eabd51c69d0c42a15d0fad/examples/chat) adapted to work with two users and increased security. The client is a minimalisticly styled one page html file.

## Security

Each chat is authenticated via a chat key that needs to be known only by the two users. If a third user obtains the key, then they can pose as one of the users. When two users participate in the chat, no new IP can join the chat.

If both users have non-malicious browsers that know how to handle TLS and the, the html file is served correctly, and the domain of the server is correctly obtained from the configured DNS server, then security of the chat is guarranteed.