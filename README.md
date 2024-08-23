## Projeto LiveStream: API

Esse repositório contém a base de código necessária para permitir que um servidor 
de livestreams seja controlado de forma organizada e eficiente. A aplicação fornece
uma interface para adição de usuários e livestreams de forma persistente, além de
conter a estrutura que recebe as transmissões e as distribui para os consumidores.

### Como executar

Para executar o projeto em sua máquina, clone o repositório e acesse a pasta contendo
o projeto. Antes de tudo, crie um arquivo `.env` seguindo o exemplo do arquivo `.env.local`.
É necessário que você tenha o `docker` instalado em sua máquina. Para 
subir todos os projetos, basta executar:

```
docker compose up --build
```

Isso irá executar os três serviços contidos no arquivo `compose.yml`

- `database`: O banco de dados da aplicação (MongoDB).
- `ls-server`: A API em Go que permite gerenciar a criação de lives e usuários.
- `nginx`: O serviço do `nginx` permite que usuários possam transmitir e consumir
streams de forma eficiente.

O servidor de ingestão se econtra na porta `8000`. A API se encontra na porta `3333`.

### Documentação da API

A documentação da API é feita utilizando o Swagger. Para acessar a documentação,
execute os seguintes comandos (o ambiente já deve estar gerado seguindo a seção
anterior).

```
# Para gerar o arquivo swagger.yaml (necessário pela primeira vez ou quando a documentação é atualizada.
docker exec ls-server make swagger

# Para lançar o servidor que disponibiliza a documentação
docker exec ls-server make swagger_serve
```

Após isso, basta acessar o endereço `http://localhost:4004/docs`

### Como faço para abrir uma live?

Para transmitir uma live, um usuário deve utilizar algum software de encoding
de streams de vídeo, como o OBS. Para transmitir, ele deve utilizar sua chave 
de stream, bem como seu nome de usuário e senha (provisionados pela API no controle de usuários),
e setar o endereço do servidor como:

```
rtmp://localhost:1936/livestream/<sua_stream_key>
```

Após a API verificar a identidade do usuário, ele já estará transmitindo dados,
mas a stream só deve começar quando ele permitir na interface web. 
