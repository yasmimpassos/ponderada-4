# Racha Histórico

Aplicativo para dividir gastos em grupo. Serve para quem viaja, divide aluguel, faz um rolê ou qualquer situação em que várias pessoas pagam coisas e no fim precisam acertar quem deve para quem. Em vez de anotar tudo no papel ou fazer conta de cabeça, o app registra cada despesa, divide automaticamente entre os participantes e mostra de forma clara o saldo de cada um.

## O que dá para fazer

1. Criar uma conta e entrar com email e senha.
2. Montar grupos para cada situação (uma viagem, a república, um churrasco).
3. Convidar gente para o grupo enviando um link de convite.
4. Registrar uma despesa dizendo quanto custou e quem participou.
5. Escanear o comprovante com a câmera para preencher o valor e a descrição sem digitar.
6. Dividir o valor automaticamente entre as pessoas escolhidas (cinquenta reais entre duas pessoas vira vinte e cinco para cada).
7. Ver dentro de cada grupo quem pagou o quê e como ficou o saldo de todo mundo.
8. Ver num lugar só todas as suas dívidas juntando todos os grupos, sabendo exatamente quem te deve e para quem você deve.
9. Marcar uma dívida como paga ou marcar que recebeu, mantendo os saldos sempre atualizados.
10. Receber um aviso na tela sempre que uma nova despesa é registrada.

## Como funciona na prática

### Entrada
Ao abrir, o app pede email e senha. Quem ainda não tem conta consegue criar uma na hora. Se você abrir o aplicativo a partir de um link de convite, depois de entrar você já cai direto no grupo para o qual foi chamado.

### Seus grupos
Logo após entrar você vê a lista de todos os grupos dos quais participa e pode criar um novo. Dessa mesma tela dá para abrir o resumo geral das suas dívidas.

### Dentro de um grupo
Cada grupo tem três visões. Uma mostra todas as despesas lançadas, com o valor e o nome de quem pagou. Outra mostra o balanço, ou seja, quanto cada pessoa do grupo deve ou tem a receber. A terceira mostra quem são os membros. Dali você também convida mais gente e adiciona uma nova despesa.

### Registrar uma despesa
Você informa a descrição, o valor e a data, e escolhe quais membros entraram naquele gasto. O valor é repartido igualmente entre os escolhidos. Se preferir, aponta a câmera para o comprovante e o app tenta preencher o valor e a descrição sozinho.

### Suas dívidas
Reúne tudo de todos os grupos numa visão única. De um lado, as pessoas que você precisa pagar. Do outro, as pessoas que te devem. Quando alguém acerta, basta confirmar o pagamento ou o recebimento que os saldos são recalculados, inclusive dentro dos grupos.

### Convite
Ao compartilhar um grupo, o app usa o compartilhamento do próprio celular, então você manda o convite por WhatsApp, email ou qualquer outro aplicativo. Quem recebe o link entra no grupo, criando conta antes caso ainda não tenha.

## Tecnologias

O aplicativo é feito em Flutter. O servidor é escrito em Go e guarda as informações em um banco de dados MySQL. A leitura dos comprovantes usa o serviço de reconhecimento de texto do Google. As notificações aparecem direto no aparelho, e o compartilhamento e a câmera usam os recursos nativos do celular.

## Como rodar

### O que você precisa
Tenha instalado o Flutter, a linguagem Go e um banco MySQL rodando.

### Banco de dados
Crie um banco e um usuário no MySQL e rode os arquivos da pasta `backend/migrations` na ordem numérica para montar as tabelas. Um exemplo de comando para cada arquivo:

```bash
mysql -u racha -p racha < backend/migrations/001_create_users.sql
```

### Servidor
Dentro da pasta `backend`, crie um arquivo `.env` com os seguintes valores:

```bash
DATABASE_URL=racha:racha@tcp(localhost:3306)/racha?parseTime=true
JWT_SECRET=umsegredoqualquer
GOOGLE_VISION_API_KEY=suachavedovision
PORT=8080
```

Depois suba o servidor:

```bash
cd backend
go run .
```

### Aplicativo
Dentro da pasta `app`, instale as dependências e rode:

```bash
cd app
flutter pub get
flutter run
```

Para testar no navegador:

```bash
flutter run -d web-server --web-port 5000
```

### Rodando em outro aparelho
Quando o app estiver aberto, toque duas vezes seguidas no botão de entrar para informar o endereço do servidor. Isso é útil quando você expõe o backend com uma ferramenta como o ngrok e quer acessar de um celular real.
