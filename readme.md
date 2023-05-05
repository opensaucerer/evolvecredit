# Zeina MFI - Evolve Credit

Zeina MFI needs software to manage deposits, savings, and withdrawals (all off the counter without payments automation). They want to be able to

- Accept deposits from their customers and keep these records
- Lock a portion of deposits on behalf of the customers for a specified duration
- Manage withdrawals
- Keep logs of all these transactions
- Transactions should be in accurate sync with customer account balances

## Running the application

- Clone the repository
- Set the `ENV_PATH` environment variable to the path of the `.env` file

```bash
export ENV_PATH=/path/to/.env
```

- Download the dependencies and run the application

```bash
make run
```

## Running the tests

- Clone the repository
- Set the `ENV_PATH` environment variable to the path of the `.env` file

```bash
export ENV_PATH=/path/to/.env
```

- Download the dependencies and run the tests

```bash
make test
```

### API Documentation

Link to [Postman Documentation](https://documenter.getpostman.com/view/11854111/2s93eVYZa8)

### Improvements

> There are no best solutions, only best trade-offs.

- The framework, barf, can be a lot more remarkable by making it into a package. BARF is an acronym for Basically, A Remarkable Framework...in reference to BARF from the marvel cinematic universe (Binarilly Augmented Retro-Framing).
- I choose to build my own framework just in case it would give me an edge ahead of others in the context of golang understanding and problem solving.
- Subrouters, GraphQL, and Testing support can be added to the framework and I will be working on this in the coming days. I love engineering things.
- A better validation layer can be implemented as a middleware. I added this as a comment somewhere in the code.
- The test suite can be made better to conver a lot more cases. Although, I'm not of the opinion of mocking databases I think unit testing some of the repository functions can be done via mocking. I didn't do this.
- Security can be improved with authentication layers (hence the reason for including `roles`) as opposed to the current use of app token in the header. Also, the current app token flow can be improved using asymmetric encryption or zero-knowledge proof.
- The account creation flow can include other account types as defined in the `AccountType` iota.
- I wasn't really sure if unlocking was meant to be automated or manual but since I kept most engineering focus on `without payment automation`, I made it manaul. However, this can be automated by adding a cron job that runs at some configured interval to unlock funds whose lock duration has expired.
- Pagination can be added to endpoints that return lists of data.
- When generating account numbers, concurrent requests can potentially generate the same account number. This can be improved by using a mutex lock or a channel to ensure that only one request is generating an account number at a time.
- I did not include foreign keys in my tables while means PostgreSQL is unaware of the relationships they share. To enforce data consistency and avoid the possibility of orphaned rows, it would be ideal to include foreign keys. 
