##### json request body for `construction/preprocess`

send this to construction/preprocess and follow the flow of operations https://www.rosetta-api.org/docs/flow.html

see https://github.com/imerkle/rosetta-solana-go/blob/master/services/construction_service_test.go#L165 for example


#### NATIVE SOL Transfer `System__Transfer`
```
{
    "network_identifier": {
        "blockchain": "solana",
        "network": "devnet"
    },
    "operations": [
        {
            "operation_identifier": {
                "index": 0
            },
            "type": "System__Transfer",
            "account": {
                "address": "HJGPMwVuqrbm7BDMeA3shLkqdHUru337fgytM7HzqTnH"
            },
            "amount": {
                "value": "-1000",
                "currency": {
                    "symbol": "SOL",
                    "decimals": 9
                }
            }
        },
        {
            "operation_identifier": {
                "index": 1
            },
            "type": "System__Transfer",
            "account": {
                "address": "CgVKbBwogjaqtGtPLkMBSkhwtkTMLVdSdHM5cWzyxT5n"
            },
            "amount": {
                "value": "1000",
                "currency": {
                    "symbol": "SOL",
                    "decimals": 9
                }
            }
        }
    ]
}

```
#### SPL TOKEN TRANSFER NEW `SplToken__TransferWithSystem`

this abstracts away the need to deal with token accounts

```
{
    "network_identifier": {
        "blockchain": "solana",
        "network": "devnet"
    },
    "operations": [
        {
            "operation_identifier": {
                "index": 0
            },
            "type": "SplToken__TransferWithSystem",
            "account": {
                "address": "HJGPMwVuqrbm7BDMeA3shLkqdHUru337fgytM7HzqTnH" // system account
            },
            "amount": {
                "value": "-1",
                "currency": {
                    "symbol": "3fJRYbtSYZo9SYhwgUBn2zjG98ASy3kuUEnZeHJXqREr",
                    "decimals": 2
                }
            },
        },
        {
            "operation_identifier": {
                "index": 1
            },
            "type": "SplToken__TransferWithSystem",
            "account": {
                "address": "CgVKbBwogjaqtGtPLkMBSkhwtkTMLVdSdHM5cWzyxT5n"  // system account
            },
            "amount": {
                "value": "1",
                "currency": {
                    "symbol": "3fJRYbtSYZo9SYhwgUBn2zjG98ASy3kuUEnZeHJXqREr",
                    "decimals": 2
                }
            },    
        }
    ]
}

```
#### SPL TOKEN TRANSFER `SplToken__Transfer`

transfer spl with token accounts

```
{
    "network_identifier": {
        "blockchain": "solana",
        "network": "devnet"
    },
    "operations": [
        {
            "operation_identifier": {
                "index": 0
            },
            "type": "SplToken__Transfer",
            "account": {
                "address": "95Dq3sXa3omVjiyxBSD6UMrzPYdmyu6CFCw5wS4rhqgV" // source token account
            },
            "amount": {
                "value": "-1",
                "currency": {
                    "symbol": "3fJRYbtSYZo9SYhwgUBn2zjG98ASy3kuUEnZeHJXqREr",
                    "decimals": 2
                }
            },
            "metadata": {
                "authority": "HJGPMwVuqrbm7BDMeA3shLkqdHUru337fgytM7HzqTnH"//required source system adress
            }
        },
        {
            "operation_identifier": {
                "index": 1
            },
            "type": "SplToken__Transfer",
            "account": {
                "address": "GyUjMMeZH3PVXp4tk5sR8LgnVaLTvCPipQ3dQY74k75L"  // token account
            },
            "amount": {
                "value": "1",
                "currency": {
                    "symbol": "3fJRYbtSYZo9SYhwgUBn2zjG98ASy3kuUEnZeHJXqREr",
                    "decimals": 2
                }
            },
            "metadata": {
                "authority": "HJGPMwVuqrbm7BDMeA3shLkqdHUru337fgytM7HzqTnH" //required source system adress
            }            
        }
    ]
}
```

#### SPL TOKEN TRANSFER NEW `SplToken__TransferNew`

```
{
    "network_identifier": {
        "blockchain": "solana",
        "network": "devnet"
    },
    "operations": [
        {
            "operation_identifier": {
                "index": 0
            },
            "type": "SplToken__TransferNew",
            "account": {
                "address": "95Dq3sXa3omVjiyxBSD6UMrzPYdmyu6CFCw5wS4rhqgV" // source token account
            },
            "amount": {
                "value": "-1",
                "currency": {
                    "symbol": "3fJRYbtSYZo9SYhwgUBn2zjG98ASy3kuUEnZeHJXqREr",
                    "decimals": 2
                }
            },
            "metadata": {
                "authority": "HJGPMwVuqrbm7BDMeA3shLkqdHUru337fgytM7HzqTnH"//required source system adress
            }
        },
        {
            "operation_identifier": {
                "index": 1
            },
            "type": "SplToken__TransferNew",
            "account": {
                "address": "CgVKbBwogjaqtGtPLkMBSkhwtkTMLVdSdHM5cWzyxT5n"  // destination system account
            },
            "amount": {
                "value": "1",
                "currency": {
                    "symbol": "3fJRYbtSYZo9SYhwgUBn2zjG98ASy3kuUEnZeHJXqREr",
                    "decimals": 2
                }
            },
            "metadata": {
                "authority": "HJGPMwVuqrbm7BDMeA3shLkqdHUru337fgytM7HzqTnH"//required source system adress
            }            
        }
    ],
    "metadata": {
        "blockhash": "42gAeAs9JE1bzqjGQtprYcdi5KyZAQeDLYVoyVSpRLTA",
        "fee_calculator": {
            "lamportsPerSignature": 5000
        }
    }
}
```
#### Spl Associated Token Account CREATE `SplAssociatedTokenAccount__Create`

Creates new spl token account for reciever
```
{
    "network_identifier": {
        "blockchain": "solana",
        "network": "devnet"
    },
    "operations": [
        {
            "operation_identifier": {
                "index": 0
            },
            "type": "SplAssociatedTokenAccount__Create",
            "account": {
                "address": "HJGPMwVuqrbm7BDMeA3shLkqdHUru337fgytM7HzqTnH" //sender/signer
            },
            "metadata": {
                "mint": "GmrqGgTJ2mmNDvqaa39NAnzcwyXtm5ntTa41zPTHyc9o" //spl token mint address
                "wallet": "42jb8c6XpQ6KXxJEHSWPeoFvyrhuiGvcCJQKumdtW78v" //reciever
            }
        },
    ]
}
```


##### json request body for `/call`


```
{
    "network_identifier": {
        "blockchain": "solana",
        "network": "devnet"
    },
    "method": "getProgramAccounts",
    "parameters": {"param": ["Feat1YXHhH6t1juaWF74WLcfv4XoNocjXA6sPWHNgAse"]}
}
```
```
{
    "network_identifier": {
        "blockchain": "solana",
        "network": "devnet"
    },
    "method": "getClusterNodes",
    "parameters": {}
}
```