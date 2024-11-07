# Mini - A Minimal Cosmos SDK Chain with Checkers Game

This repository contains a working Cosmos SDK chain that implements a checkers game module. It uses minimal modules and serves as a starting point for building your own chain with game functionality.

## Overview

`Minid` is built with the latest version of [Cosmos-SDK](https://github.com/cosmos/cosmos-sdk) and includes:
- Basic chain functionality
- Checkers game module
- Account management
- Token handling (mini tokens)

## Prerequisites

* Go 1.21 or later ([installation guide](https://go.dev/doc/install))
* Configure your environment:
  ```bash
  export PATH="$PATH:/usr/local/go/bin:$(/usr/local/go/bin/go env GOPATH)/bin"
  ```

## Installation & Setup

1. Clone and build the chain:
```bash
git clone git@github.com:cosmosregistry/chain-minimal.git
cd chain-minimal
make install   # Install the minid binary
make init      # Initialize the chain
```

2. Start the chain:
```bash
minid start --minimum-gas-prices="0mini"
```

## Playing Checkers

### Create a New Game

```bash
minid tx checkers create-game <game-id> \
    <black-player-address> \
    <red-player-address> \
    --from <creator-address> \
    --gas auto \
    --gas-prices 0mini \
    --chain-id demo \
    --yes
```

Example:
```bash
minid tx checkers create-game myGame1 \
    mini1dyr9fktej5af9mq3vdwv76uy9dvptkaxqnqm24 \
    mini10gl9v6utc7hz6dflasdaw3d29z9368ynphh683 \
    --from alice --yes
```

### Query Game Status

```bash
minid query checkers get-game <game-id>
```

Example:
```bash
minid query checkers get-game myGame1
```

Example output:
```
Game:
  black: mini1dyr9fktej5af9mq3vdwv76uy9dvptkaxqnqm24
  board: '*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*'
  red: mini10gl9v6utc7hz6dflasdaw3d29z9368ynphh683
  turn: b
```

### Export Chain State

To export the current state of the chain:
```bash
minid export
```

This will output the entire chain state including:
- Account balances
- Game states
- Validator information
- Module parameters

## Key Management

After initialization, two accounts are created:
1. Alice (with 10,000,000 mini tokens)
2. Bob (with 1,000 mini tokens)

Keep the generated mnemonics safe - they are required to recover these accounts.

## Troubleshooting

1. If `minid` is not found after installation:
   - Verify your `$PATH` configuration
   - Run `which minid` to check the installation
   - Try re-running `make install`

2. If you get a gas price error:
   - Always include `--gas-prices` in your transactions
   - Or set minimum gas prices in ~/.minid/config/app.toml

3. If a game query returns empty:
   - Verify the game ID is correct
   - Ensure the transaction to create the game was successful
   - Check the chain is running and synced

## Development

The chain includes several key components:
- Custom checkers module
- Cosmos SDK core modules (auth, bank, staking)
- Minimal dependency set for optimal performance

## Useful Links

* [Cosmos-SDK Documentation](https://docs.cosmos.network/)
* [Cosmos Network](https://cosmos.network/)
* [Tendermint Core](https://tendermint.com/)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the terms of the MIT license.