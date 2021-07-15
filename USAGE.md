##### json request body for `construction/payloads`

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
    ],
    "metadata": {
        "blockhash": "42gAeAs9JE1bzqjGQtprYcdi5KyZAQeDLYVoyVSpRLTA",
        "fee_calculator": {
            "lamportsPerSignature": 5000
        }
    }
}
```
#### SPL TOKEN TRANSFER `SplToken__Transfer`

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
                "address": "95Dq3sXa3omVjiyxBSD6UMrzPYdmyu6CFCw5wS4rhqgV" // token account
            },
            "amount": {
                "value": "-1",
                "currency": {
                    "symbol": "3fJRYbtSYZo9SYhwgUBn2zjG98ASy3kuUEnZeHJXqREr",
                    "decimals": 2
                }
            },
            "metadata": {
                "authority": "HJGPMwVuqrbm7BDMeA3shLkqdHUru337fgytM7HzqTnH"
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
                "authority": "HJGPMwVuqrbm7BDMeA3shLkqdHUru337fgytM7HzqTnH"
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
                "address": "95Dq3sXa3omVjiyxBSD6UMrzPYdmyu6CFCw5wS4rhqgV" // token account
            },
            "amount": {
                "value": "-1",
                "currency": {
                    "symbol": "3fJRYbtSYZo9SYhwgUBn2zjG98ASy3kuUEnZeHJXqREr",
                    "decimals": 2
                }
            },
            "metadata": {
                "authority": "HJGPMwVuqrbm7BDMeA3shLkqdHUru337fgytM7HzqTnH"
            }
        },
        {
            "operation_identifier": {
                "index": 1
            },
            "type": "SplToken__TransferNew",
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
            "metadata": {
                "authority": "HJGPMwVuqrbm7BDMeA3shLkqdHUru337fgytM7HzqTnH"
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
    ],
    "metadata": {
        "blockhash": "42gAeAs9JE1bzqjGQtprYcdi5KyZAQeDLYVoyVSpRLTA",
        "fee_calculator": {
            "lamportsPerSignature": 5000
        }
    }
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