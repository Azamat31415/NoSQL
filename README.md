# API Documentation
## 1) Products

| Operation | HTTP Method | URL                         | Expected Status |
|-----------|-------------|-----------------------------|-----------------|
| Create    | POST        | /products                   | 201 Created     |
| Read      | GET         | /products/{id} or /products | 200 OK          |
| Update    | PUT         | /products/{id}              | 200 OK          |
| Delete    | DELETE      | /products/{id}              | 204 No Content  |
 
Request Body:
    {
        "name": "Dog's food",
        "description": "Very nice food",
        "price": 30.99,
        "stock": 10,
        "category": "dogs",
        "subcategory": "feed",
        "type": "dry"
    }

## 2) Authentication
### 2.1) Register
 
| Operation | HTTP Method | URL       | Expected Status |
|-----------|-------------|-----------|-----------------|
| Register  | POST        | /register | 201 Created     |

Request Body:
    {
        "email": "azabraza061005@gmail.com",
        "password": "aza061005",
        "first_name": "Baiken",
        "last_name": "Ashimov",
        "address": "Astana",
        "phone": "+77072050305"
    }

### 2.2) Login

 | Operation | HTTP Method | URL    | Expected Status |
 |-----------|-------------|--------|-----------------|
 | Login     | POST        | /login | 200 OK          |
 
Request Body:
    {
        "email": "azabraza061005@gmail.com",
        "password": "aza061005"
    }

## 3) Subscription
 
| Operation | HTTP Method | URL                 | Expected Status |
|-----------|-------------|---------------------|-----------------|
| Create    | POST        | /subscriptions      | 201 Created     |
| Delete    | DELETE      | /subscriptions/{id} | 204 No Content  |

Request Body:
    {
        "user_id": 1,
        "interval_days": 30,
        "type": "premium",
        "status": "active"
    }

## 4) Orders
 
| Operation               | HTTP Method | URL                                 | Expected Status   |
|-------------------------|-------------|-------------------------------------|-------------------|
| Create Order            | POST        | /orders                             | 201 Created       |
| Update Order Status     | PUT         | /orders/{id}/status/update          | 200 OK            |
| Choose Delivery Method  | PUT         | /orders/{order_id}/delivery         | 200 OK            |
| Get All Orders          | GET         | /orders                             | 200 OK            |
| Get Order History       | GET         | /order-history/{user_id}            | 200 OK            |
    
Request Body:
    i) Create order
        {
            "user_id": 1,
            "delivery_method": "courier",
            "address": "123 Main St, City, Country",
            "total_price": 150.75,
            "order_items": [
                {
                    "product_id": 101,
                    "quantity": 2,
                    "price": 50.00
                },
                {
                    "product_id": 102,
                    "quantity": 1,
                    "price": 50.75
                }
            ]
        }
    ii) Update order status
        {
            "status": "shipped"
        }
    iii) Choose delivery method
        {
            "delivery_method": "pickup",
            "pickup_point_id": 5
        }

## 5) User's Address

| Operation          | HTTP Method | URL                 | Expected Status |
|--------------------|-------------|---------------------|-----------------|
| Get User's Address | GET         | /users/{id}/address | 200 OK          |

## 6) Cart

| Operation                   | HTTP Method | URL                                      | Expected Status   |
|-----------------------------|-------------|------------------------------------------|-------------------|
| Add to Cart                 | POST        | /cart                                    | 201 Created       |
| Remove from Cart            | DELETE      | /cart/{id}                               | 204 No Content    |
| Update Cart Item Quantity   | PUT         | /cart/update/{id}/{quantity}             | 200 OK            |
| Remove One Item from Cart   | DELETE      | /cart/{id}/byone                         | 200 OK            |
| Get Cart by User            | GET         | /cart/user/{user_id}/products            | 200 OK            |
| Get Cart ID                 | GET         | /cart/{user_id}/{product_id}             | 200 OK            |

Request Body: 
    i) Add to cart
        {
            "user_id": 1,
            "product_id": 101,
            "quantity": 3
        }
    ii) Update cart item quantity
        {
            "quantity": 5
        }
    iii) Remove one item from cart
        {
            "user_id": 1,
            "product_id": 101
        }

# Document database schema
## Product
{
    "ID": 1,
    "CreatedAt": "2025-01-26T15:40:46.44256+05:00",
    "UpdatedAt": "2025-01-28T17:59:10.842759+05:00",
    "DeletedAt": null,
    "name": "Dog food",
    "description": "Delicious dog food",
    "price": 20.99,
    "stock": 10,
    "category": "Dog",
    "subcategory": "",
    "type": ""
}

## Orders and Order History
{
    "ID": 12,
    "UserID": 1,
    "DeliveryMethod": "courier",
    "PickupPointID": null,
    "Address": "123 Main St",
    "Status": "pending",
    "TotalPrice": 99.99,
    "CreatedAt": "2025-02-05T16:16:19.549169+05:00",
    "UpdatedAt": "2025-02-05T16:16:19.549169+05:00",
    "OrderItems": []
}

## User's Address
{
    "address": "Astana"
}
