# pdash

orders dashboard in microservices architecture

![image](https://user-images.githubusercontent.com/48713070/187465539-cd64ef2c-9876-49a1-b89e-dd9d8062562e.png)

## run with

```sh
docker-compose up --build -d
```

## urls

- dashboard: [http://localhost:3000](http://localhost:3000)
- customers service: [http://localhost:8001/customers](http://localhost:8001/customers)
- customers swagger: [http://localhost:8001/swagger/](http://localhost:8001/swagger/)
- orders service: [http://localhost:8002/orders](http://localhost:8002/orders)
- orders swagger: [http://localhost:8002/swagger/](http://localhost:8002/swagger/)
- suppliers service: [http://localhost:8003/suppliers](http://localhost:8003/suppliers)
- suppliers swagger: [http://localhost:8003/swagger/](http://localhost:8003/swagger/)
- auth service: [http://localhost:8004/users](http://localhost:8004/users)
- auth swagger: [http://localhost:8004/swagger/](http://localhost:8004/swagger/)

## stop with

```sh
docker-compose down
```

## versions

- docker: `20.10.17`
- docker-compose: `1.29.2`
