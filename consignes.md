# Mid Backend Exercice

# Exercise: Tezos Delegation Service

In this exercise, you will build a Golang service that gathers new [delegations](https://opentezos.com/node-baking/baking/delegating/) made on the Tezos protocol and exposes them through a public API. 

## Requirements:

The solution is composed of two parts:

- It must poll delegations:
    - It must continuously poll new delegations from this tzkt API endpoint: https://api.tzkt.io/#operation/Operations_GetDelegations
    - For each delegation, save the following information: sender's address, timestamp, amount, and block height.
    - The data aggregation must store the delegations data in a persistent store of your choice.
    - Indexing historical data is a bonus.
- It must expose the collected data through a public API endpoint:
    - The endpoint must be available at: `GET /xtz/delegations`
    - The API must read data from the store.
    - The response format must be:
    
    ```jsx
    {
      "data": [ 
        {
            "timestamp": "2022-05-05T06:29:14Z",
            "amount": "125896",
            "delegator": "tz1a1SAaXRt9yoGMx29rh9FsBF4UzmvojdTL",
            "level": "2338084"
        },
        {
            "timestamp": "2021-05-07T14:48:07Z",
            "amount": "9856354",
            "delegator": "KT1JejNYjmQYh8yw95u5kfQDRuxJcaUPjUnf",
            "level": "1461334"
        }
      ],
    }
    ```
    
- The sender’s address is the delegator.
- The delegations must be listed most recent first.
- The endpoint takes one optional query parameter `year`, which is specified in the format YYYY and will result in the data being filtered for that year only.

### Additional notes

- The code must be tested.
- How to run the solution locally must be simple and documented.
- The solution must thrive to be simple while fulfilling all the requirements.

Please share a archive ( `zip` , `tar` or equivalent) of your git project via email.