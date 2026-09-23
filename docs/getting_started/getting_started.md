## Project Setup

### Set environment variables

Define the environment variables in their relevant folders according to relevant their .env.example files.

### Up all services

Run below command from the main root of the project. It will up all ther services except the crawler.

```bash
docker compose up --build
```
Crawler works as a corn job. Use below command to run the crawler.

```bash
docker compose --profile crawler run --rm crawler
```

### ThunderID Setup

Once all services are running, access the ThunderID console at "https://thunder:8090/console". and log in using the credentials provided in your .env file to open the main dashboard of ThunderID. Click the Import button, then drag and drop the [thunderid-environment.env](./thunderid_configs/thunderid-environment.env) and [thunderid-config.yml](./thunderid_configs/thunderid-config.yml) files. Finally, generate the client secrets for both applications and add them to the relevant fields in your .env file.

For more details :- [ThunderID Documentaion](https://thunderid.dev/docs/)

### MCP Setup

Once all services are running, access the MCP at "http://localhost:9090/mcp". MCP server is also connected with ThunderID and This server is using DCR (Dynamic Client Registration) method. Refer documentation to [setup MCP server with ThunderID](https://thunderid.dev/docs/getting-started/connect-your-mcp/python/).