# Authorizer
Autorização de transações financeiras.

## Como funciona?
- Hey, senhor, posso seguir com está compra?
- Vamo ve. Se não tiver nenhuma violão de regra, ta liberado.

## E como eu uso?

* Build: `docker build -t authorizer .`
* Run: `docker run -it authorizer`

Dentro do console do container, basta aplicar o stream do arquivo ao binário compilado.

* `./authorizer < operations`

### Dependências
Nenhuma foi utilizada nesse projeto o/
(não que libs sejam ruins, libs eu amo vocês)

# Arquitetura e decisões
Como base, foi utilizado os conceitos de **Clean Architecture** para criar o projeto,com algumas modificações mirando simplicidade da aplicação.

Como estratégia de dados, utilizei **Event Sourcing** para gerenciar o estado da aplicação.

## Por que Event Sourcing?
* Imutabilidade de dados, aonde não há atualização de dados em si, mas somente a sumarização dos eventos ocorridos para chegar no estado atual.
* Facilidade para entender como a aplicação estava no passado e com isso tomar decisões. Temos um auditlog completo e confiável.
* Facilidade para desacoplar o código.
* Facilidade de extender o código para outros casos de uso.

Tudo isso, unido com a natureza da aplicação, que autoriza transações baseadas em o quê ocorreu no passado, **Event Sourcing** se mostrou uma boa solução.

## E quais os contras?
* Maior complexidade para ter o estado final da aplicação.
* Maior complexidade para gerenciar os dados. Não usei aggregates para gerar gerenciar o estado e os eventos para não adicionar mais essa complexididade. Centralizei nos Repository que já tinham a responsabilidade de persistir os dados.


## CQRS
Não há necessidade neste cenário de utilizar CQRS, que é comum se ver utiliando com **Event Sourcing**. Além disso, é importante ter o estado mais acurado possível para tomar a decisão se a transação deve ser autorizada ou não.

Mesmo não utiliando agora, a aplicação permite facilmente extender para usar CQRS e replicar os dados para outros uso. Para isso, basta implementar um **Repository** que replique esses eventos através de uma mensageria.

## Paradgimas
Em geral, a aplicação é Orientada a Objetos e usa alguns benefícios de programação funcional. Um exemplo disso, são as checagens de violações, que são funções puras que executam em cima de uma estrutura que contém os dados necessários para validar.

Por ser funções puras e usar dados imutáveis que são passados como cópias, elas podem ser rodadas em paralelo e assim acelar as checagens (mas ficou por uma v2 devido ao tempo. DESCULPA :///).

## Por que Golang?
Não é a linguagem que mais tenho experiência, sendo elas JavaScript e Ruby. Porém, como tive experiências passadas com Golang, utilizei para ganhar tempo e aprender mais. Rich Hickey que me perdoe, mas não me senti seguro para criar uma aplicação com tempo fechado em Clojure. (mas deu vontade)

## Testes
Os testes de Casos de Uso (Use Cases) são mais de integração e o restante mais unitários. Não foi coberto 100%, mas as funcionalidades base estão cobertas.

## Melhorias futuras
* Adicionar camadas de abstração para desacoplar mais o código e facilitar trocar. Como adicionar Controllers e Services (ficou método de **service** dentro de **repository**).
* Melhorias nas interfaces da aplicação para facilitar trocas de estruturas.
* Maior uso de go routines e channels para usar o poder do golang e rodar tudo em paralelo. Lei de Moore me desculpa.
* Estratégias para evitar percorrer todos os eventos sempre que precisar de um estado, como guardar a referência do último evento sumarizado ou outras estratégias.
* Rodar a aplicação passando um stream de eventos fora do container Docker (hoje tem que ta dentro).
* Melhoria nos comentários com maiores informações e exemplos de uso. Priorizei a funcionalidade x tempo.
* Com certeza tem mais coisas que ainda não vi ainda :P

## Processo de criação
* Desenhei mais ou menos como queria a arquitetura do projeto em um caderno, o que deu umas 3 folhas de rascunho depois de pensar em algumas versões.
* Fiz uma POC simples para testar o storage em memória e ter mais familiaridade com a linguagem
* Quebrei o projeto em algumas tarefas para ter um guia e mais ou menos o tempo que ia levar
* Não usei um fluxo de Kanban em si, mas fui dando check nas camadas de abstração que fui criando
* E ai foi só codar

No mais, é isso galera! Curti MUITO e aprendi MUITO, valeu!

ps: quase que chamei o projeto de Julius em homenagem a Everybody Hates Chris.