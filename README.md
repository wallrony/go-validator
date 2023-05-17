# Validação de DTO

A validação de DTO tem um processo mais automatizado e simples, para evitar duplicação de código e para facilitar o processo de identificação de informações incorretas. Para isso, utilizamos um `Validator Helper`, que auxilia a identificar os erros via tags de atributos nos DTOs.

## Problema

Primeiro, vamos assumir que temos um DTO chamado `Account`, considerando a seguinte estrutura:

```go
type Account struct {
	Name      string `json:"name"`
	BirthDate string `json:"birth_date"`
	Email     string `json:"email"`
}
```

Em cada atributo informamos que há uma tag chamada `json` cujo valor é um nome do atributo no padrão `snake_case`. Quando temos um conjunto de dados (ex.: um Map/JSON) que será encaixado nessa estrutura, além de ser colocado cada valor no seu devido lugar, temos que também realizar um processo de validação de valores, como no nome, email, data de nascimento e nos demais atributos. De forma manual, adicionaríamos várias linhas de código para podermos validar cada atributo (considerando também que teríamos vários lugares do código utilizando a mesma lógica ou outras estruturas com atributos semelhantes).

Nesse caso, teríamos o seguinte código considerando a estrutura fornecida:

```go
func main() {
	timeLayout := "2006-01-02"
	data := map[string]interface{}{
		"name":       "Test Man",
		"birth_date": "2000-01-01",
	}
	serialized := request.Account{}
	serialized.Name = data["name"].(string) // might fail
	serialized.BirthDate = data["birth_date"].(string) // might fail
	if serialized.Name == "" {
		log.Fatal("Name is required!")
	} else if len(strings.Split(serialized.Name, ",")) == 1 {
		log.Fatal("You need to provide the first and the last name!")
	} else if v, err := time.Parse(timeLayout, serialized.BirthDate); err != nil {
		log.Fatal("Invalid date format!")
	} else if v.After(time.Now()) {
		log.Fatal("Invalid birth date!")
	}
}
```

No método acima, estamos verificando alguns valores atribuídos na variável `serialized`. Porém com o tempo teremos mais estruturas com atributos parecidos e vários lugares utilizando da mesma validação. Além disso, temos que pensar na tipagem the cada valor passado no map `data` (como colocado quando atribuindo os atributos `Name` e `BirthDate`).

## Solução

O `Validator Helper` foi criado para cobrir esses casos faltantes e conseguirmos ter uma menor quantidade de linhas para resolver problemas de validação de conjunto de dados. Considerando cada atributo, por padrão, como opcional, é possível informar quais atributos são obrigatórios e/ou têm validações específicas.

Para utilizarmos, precisamos utilizar o seguinte trecho de código:

```go
...
	dto, err := validator.ValidateDTO[Account](data)
	fmt.Println("DTO:", dto)
	fmt.Println("ERROR:", err)
...
```

O método `ValidateDTO` vai utilizar um tipo genêrico para basear a conversão de valores presentes dentro dos dados da variável `data`. Essa função retorna o dto convertido e um erro do tipo `ValidationError`. Caso os dados não sejam válidos e existam erros relacionados, a variável `dto` será `nil` e terá um erro. Caso contrário, a variável `dto` irá conter os valores esperados. Neste caso, temos nome e data de nascimento:

```bash
# go run main.go
DTO: &{Test Man 2000-01-01 }
ERROR: <nil>
```

Nesse ponto, não informamos nenhuma validação específica e o validador conseguiu definir os atributos `Name` e `BirthDate`. Porém, se informarmos valores errados considerando a tipagem da estrutura, o seguinte irá acontecer:

```go
...
	data := map[string]interface{}{
		"name":       "Test Man",
		"birth_date": 2,
	}
...
```

Executando:
```bash
# go run main.go
DTO: &{Test Man  }
ERROR: <nil>
```

O valor da data de nascimento não foi definido no atributo `BirthDate` e nenhum erro foi retornado. Isso aconteceu pois a variável é facultativa e nenhuma validação foi requisitada. Para consertarmos isso e definirmos a data de nascimento como obrigatória, precisamos adicionar uma tag chamada `validate` de validação na estrutura:

```go
type Account struct {
	Name      string `json:"name"`
	BirthDate string `json:"birth_date" validate:"required"`
	Email     string `json:"email"`
}
```

Com isso, ao executar o mesmo código, o atributo `BirthDate` será considerado obrigatório e um erro será adquirido:

```bash
# go run main.go
DTO: <nil>
ERROR: 'birth_date' field type must be 'string'
```

Caso você queira que o mesmo atributo ainda seja opcional, mas quando informado (seja diferente de `nil`) seja convertido devidamente pelo tipo informado, basta adicionar a validação de tipo na tag `validate`:

```go
type Account struct {
	Name      string `json:"name"`
	BirthDate string `json:"birth_date" validate:"type"`
	Email     string `json:"email"`
}
```

Ao executar o mesmo código, teremos a seguinte saída:

```bash
# go run main.go
DTO: <nil>
ERROR: 'birth_date' field type must be 'string'
```

## Validação Parcial

Até então vimos que é possível ter uma validação bruta, onde caso der certo, teremos os dados. Caso contrário, teremos o erro. Em alguns casos talvez seja necessário ter ambos, reaproveitando os valores que a validação foi feita e o valor se encontra correto. Para isso, basta utilizarmos o método `ValidateDTOPartially`, onde a validação, por mais que nos retorne um erro, o DTO será retornado com os dados validados até então.

Para testarmos isso, considerando que temos a seguinte estrutura e dados:

```go
type Account struct {
	Name      string `json:"name" validate:"required"`
	BirthDate string `json:"birth_date" validate:"type"`
	Email     string `json:"email"`
}
...
	data := map[string]interface{}{
		"name":       "",
		"birth_date": "01/01/2000",
	}
...
```

A saída esperada será:

```bash
# go run main.go
DTO: &{ 01/01/2000 }
ERROR: 'name' field of type 'string' is missing or empty
```

## Validação de Tipos Específicos

### Data

No exemplo anterior vimos que é possível adicionar verificações obrigatórias ou não que validam sua presença ou seu tipo quando valores de um atributo específico estão presentes. Porém, dependendo do atributo, necessitamos de uma validação a mais. Um desses casos é o caso do atributo `BirthDate` que precisa de uma validação de data.

Para isso, podemos alterar a nossa estrutura para ter o seguinte valor na tag `validate`:

```go
type Account struct {
	Name      string `json:"name"`
	BirthDate string `json:"birth_date" validate:"type,date"`
	Email     string `json:"email"`
}
```

Com isso, utilizando o seguinte map de dados:

```go
...
	data := map[string]interface{}{
		"name":       "Test Man",
		"birth_date": "2000",
	}
...
```

Essa é a saída esperada:

```bash
# go run main.go
DTO: <nil>
ERROR: 'birth_date' field doesn't match with the '2006-01-02' format
```

Por padrão, o `Validator Helper` utiliza o formato `2006-01-02` para validar datas. Quando informando uma data correta (ex.: "2000-01-01"), a saída esperada é:

```bash
# go run main.go
DTO: &{Test Man 2000-01-01 }
ERROR: <nil>
```

Se você deseja definir um formato específico de validação de data, basta informar na tag `validate` seguindo a sintaxe "date=&lt;format&gt;", onde "&lt;format&gt;" pode ser "2006" caso queira somente o ano, "2006-01" caso queira o ano e o mês ou não definir nenhum formato (definindo somente "date") para pegar toda a formatação de data. Confira a saída esperada quando na tag `validate` está definido "type,date=2006-01" e passamos a data completa ("2000-01-01"):

```bash
# go run main.go
DTO: <nil>
ERROR: 'birth_date' field doesn't match with the '2006-01' format
```

E agora quando informamos a data no formato correto ("2000-01"):

```bash
# go run main.go
DTO: &{Test Man 2000-01 }
ERROR: <nil>
```

Essa definição de formatação livre de data permite com que seja possível, por exemplo, informar datas no padrão brasileiro como "01/01/2000" definindo na tag `validate` a sintaxe "date=02/01/2006":

```bash
# go run main.go
DTO: <nil>
ERROR: 'birth_date' field doesn't match with the '02/01/2006' format
```

### Email

Para validação de valores do tipo `email`, basta definir a regra `email` dentro da tag `validate` da seguinte forma:

```go
type Account struct {
	Name      string `json:"name"`
	BirthDate string `json:"birth_date" validate:"type,date=02/01/2006"`
	Email     string `json:"email" validate:"email"`
}
```

Com isso, caso informemos o seguinte valor em map:

```go
...
	data := map[string]interface{}{
		"name":       "Test Man",
		"birth_date": "01/01/2000",
		"email":      "Test",
	}
...
```

A saída esperada será:

```bash
# go run main.go
DTO: <nil>
ERROR: the value provided for the 'email' field isn't a valid email
```

### Quantidade de Caracteres

Podemos validar a quantidade de caracteres pelas regras `minlen=<something>`, `len=<something>` e `maxlen=<something>`, onde no lugar de `<something>` podemos colocar um valor que pode determinar o mínimo, específico ou máximo tamanho de um texto. Ambas as regras "minlen" e "maxlen" são inclusivas (`minlen` funciona semanticamente como "a partir de X" e `maxlen` funciona semanticamente como "até X").

> ## minlen
>
> Para utilizar a largura mínima, iremos adicionar a regra `minlen` no atributo `Name`, ficando da seguinte forma:
> 
> ```go
> type Account struct {
> 	Name string `json:"name" validate:"minlen=5"`
> }
> ```
> 
> Ao fornecer os seguintes dados:
> 
> ```go
> ...
> 	data := map[string]interface{}{
> 		"name": "Test",
> 	}
> ...
> ```
> 
> Essa será a saída esperada:
> 
> ```bash
> # go run main.go
> DTO: <nil>
> ERROR: 'name' field must have at least 5 characters
> ```
> #

> ## len
>
> Para utilizar a largura específica, iremos adicionar a regra `len` no atributo `Name`, ficando da seguinte forma:
> 
> ```go
> type Account struct {
> 	Name string `json:"name" validate:"len=5"`
> }
> ```
> 
> Ao fornecer os seguintes dados:
> 
> ```go
> ...
> 	data := map[string]interface{}{
> 		"name": "Test",
> 	}
> ...
> ```
> 
> Essa será a saída esperada:
> 
> ```bash
> # go run main.go
> DTO: <nil>
> ERROR: 'name' field must have 5 characters
> ```
> #

> ## maxlen
>
> Para utilizar a largura máxima, iremos adicionar a regra `maxlen` no atributo `Name`, ficando da seguinte forma:
> 
> ```go
> type Account struct {
> 	Name string `json:"name" validate:"len=10"`
> }
> ```
> 
> Ao fornecer os seguintes dados:
> 
> ```go
> ...
> 	data := map[string]interface{}{
> 		"name": "My Awesome and Beautiful Name",
> 	}
> ...
> ```
> 
> Essa será a saída esperada:
> 
> ```bash
> # go run main.go
> DTO: <nil>
> ERROR: 'name' field must have 5 characters at max
> ```
> #

### Slices

Para validar slices, você pode utilizar as regras `slice:len=X`, `slice:minlen=X` e `slice:maxlen=X`, onde `X` é a quantidade desejada.

Exemplo: quando é necessário cadastrar uma lista de usuários e necessitamos da lista do nome das pessoas, nós normalmente utilizaríamos a seguinte estrutura:

```go
type Account struct {
	Names []string `json:"names"`
}
```

E os seguintes valores:

```go
...
	data := map[string]interface{}{
		"names": []string{
			"test",
			"test",
		},
	}
...
```

Ao executarmos sem nenhuma regra de Slice, o resultado esperado será:

```bash
# go run main.go
DTO: &{[test test]}
ERROR: <nil>
```

Com isso, podemos explorar as regras de validação de quantidade de elementos (veja abaixo).

> ## len
>
> Para utilizar a regra de quantidade específica, iremos adicionar a regra `slice:len` no atributo `Names`, ficando da seguinte forma:
> 
> ```go
> type Accounts struct {
> 	Names []string `json:"name" validate:"slice:len=2"`
> }
> ```
> 
> Ao fornecer os seguintes dados:
> 
> ```go
> ...
>	data := map[string]interface{}{
>		"names": []string{
>			"Test",
>		},
>	}
> ...
> ```
> 
> Essa será a saída esperada:
> 
> ```bash
> # go run main.go
> DTO: <nil>
> ERROR: the 'names' field must have 2 elements
> ```
> #

> ## minlen
>
> Para utilizar a regra de quantidade mínima, iremos adicionar a regra `slice:minlen` no atributo `Names`, ficando da seguinte forma:
> 
> ```go
> type Accounts struct {
> 	Names []string `json:"name" validate:"slice:minlen=2"`
> }
> ```
> 
> Ao fornecer os seguintes dados:
> 
> ```go
> ...
>	data := map[string]interface{}{
>		"names": []string{
>			"Test",
>		},
>	}
> ...
> ```
> 
> Essa será a saída esperada:
> 
> ```bash
> # go run main.go
> DTO: <nil>
> ERROR: the 'names' field must have at least 2 elements
> ```
> #

> ## maxlen
>
> Para utilizar a regra de quantidade máxima, iremos adicionar a regra `slice:maxlen` no atributo `Names`, ficando da seguinte forma:
> 
> ```go
> type Accounts struct {
> 	Names []string `json:"name" validate:"slice:maxlen=2"`
> }
> ```
> 
> Ao fornecer os seguintes dados:
> 
> ```go
> ...
>	data := map[string]interface{}{
>		"names": []string{
>			"Test1",
>			"Test2",
>			"Test3",
>		},
>	}
> ...
> ```
> 
> Essa será a saída esperada:
> 
> ```bash
> # go run main.go
> DTO: <nil>
> ERROR: the 'names' field must have 2 elements at max
> ```
> #

### Validações de elementos

Caso queira colocar uma validação específica em cada elemento da lista, é possível somente informando o tipo de validação no atributo, assim todas as validações que não sejam de lista serão aplicadas em cada elemento.

Exemplo: ao invés de nomes, vamos pensar que é necessário adquirir uma lista de emails onde precisamos validar se cada item fornecido é um email válido. Nesse caso, teremos a seguinte estrutura:

```go
type Accounts struct {
	Emails []string `json:"emails" validate:"email"`
}
```

Considerando que estamos fornecendo a seguinte lista de emails:

```go
...
	data := map[string]interface{}{
		"emails": []string{
			"Test1",
			"test@email.com",
			"Test3",
		},
	}
...
```

Essa é a resposta esperada:

```bash
# go run main.go
DTO: <nil>
ERROR: the value provided for the 'emails[0]' field isn't a valid email & the value provided for the 'emails[2]' field isn't a valid email
```

### Objetos Aninhados

Também é possível fazer a validação de objetos aninhados. Supondo um exemplo onde temos que passar todas as informações de perfil de uma conta, podemos utilizar a seguinte estrutura:

```go
type Profile struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName"`
	Email     string `json:"email" validate:"required,email"`
}

type Account struct {
	Profile Profile `json:"profile"`
}
```

Nessa, existem regras de validação somente nos atributos dos atributos da estrutura `Profile` que está aninhada na estrutura `Account`.

Podemos utilizar o seguinte código para validar a estrutura:

```go
func main() {
	data := map[string]interface{}{}
	dto, err := validator.ValidateDTO[Account](data)
	fmt.Println("DTO:", dto)
	fmt.Println("ERROR:", err)
}
```

Essa será a saída esperada:

```bash
# go run main.go
DTO: <nil>
ERROR: 'profile.firstName' field of type 'string' is missing or empty & 'profile.email' field of type 'string' is missing or empty
```

Quando informado um mapa de valores corretos como esse:

```go
...
	data := map[string]interface{}{
		"profile": map[string]interface{}{
			"firstName": "John",
			"email":     "test@email.com",
		},
	}
...
```

Essa será a resposta esperada:

```bash
# go run main.go
DTO: &{{John  test@email.com}}
ERROR: <nil>
```

Se caso você necessite usar a mesma estrutura de dados (ex.: perfil), porém dessas, ainda de forma aninhada, você precise validar somente alguns em específico, você precisa adicionar a regra `nestedProps` da seguinte forma:

```go
type UpdateAccount struct {
	Profile Profile `json:"profile" validate:"nestedProps=firstName"`
}
```

Nesse caso temos a estrutura `UpdateAccount` que, por meio da regra `nestedProps` diz que somente o atributo `firstName`, dentre todos que têm validação, deve ser validado. Utilizando o seguinte conjunto de dados:

```go
...
	data := map[string]interface{}{
		"profile": map[string]interface{}{},
	}
...
```

Essa é a resposta esperada:

```bash
# go run main.go
DTO: <nil>
ERROR: 'profile.firstName' field of type 'string' is missing or empty
```

## Limitações (Problemas Conhecidos)

1. Não é possível validar devidamente listas de objetos aninhados.
	Exemplo: 
	```go
	type Profile struct {
		FirstName string `json:"firstName" validate:"required"`
		LastName  string `json:"lastName"`
		Email     string `json:"email" validate:"required,email"`
	}

	type UpdateAccount struct {
		Profiles []Profile `json:"profiles" validate:"required"`
	}
	```
