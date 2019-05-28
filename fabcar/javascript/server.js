var bodyParser = require('body-parser');
var express = require('/data/npm/lib/node_modules/express');
var app = express();

app.use(bodyParser.urlencoded({extended: false}));
app.use(express.static('frontEnd'));
app.use(bodyParser.json());

app.get('/', function(req, res){
	res.sendFile(__dirname + '/frontEnd/main.html');
	});

app.get('/ajax/login', function(req, res){
	var username = req.query.username;
	var password = req.query.password;
	console.log(username);
	console.log(password);

	res.send(["登陆成功！"]);
	});

app.get('/ajax/donate', function(req, res){
	var amount = req.query.amount;
	console.log(amount.toString());

	mainCreate(amount.toString());
	res.send(["捐款成功！"]);
	});
var result;
app.get('/ajax/query', async (req, res) => {
	console.log(req.query.search);

	var result = await mainQuery();
	res.send([result]);
//	result.then(res.send([result]));

/*	var result = mainQuery();
	
	console.log(result);

	res.send([result]);

	var result;
	async function test(){
		this.result = await mainQuery();
		console.log(result);
	}
*/
	});
var server = app.listen(8080);

//查询
'use strict';

const { FileSystemWallet, Gateway } = require('fabric-network');
const path = require('path');

const ccpPath = path.resolve(__dirname, '..', '..', 'first-network', 'connection-org1.json');

async function mainQuery() {
    try {

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = new FileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const userExists = await wallet.exists('user1');
        if (!userExists) {
            console.log('An identity for the user "user1" does not exist in the wallet');
            console.log('Run the registerUser.js application before retrying');
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccpPath, { wallet, identity: 'user1', discovery: { enabled: true, asLocalhost: true } });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork('mychannel');

        // Get the contract from the network.
        const contract = network.getContract('fabcar');

        // Evaluate the specified transaction.
        // queryCar transaction - requires 1 argument, ex: ('queryCar', 'CAR4')
        // queryAllCars transaction - requires no arguments, ex: ('queryAllCars')
        const result = await contract.evaluateTransaction('queryAllDeals');
        console.log(`Transaction has been evaluated, result is: ${result.toString()}`);

	return result.toString();

    } catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
        process.exit(1);
    }
}


//插入
async function mainCreate(amount) {
    try {

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = new FileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const userExists = await wallet.exists('user1');
        if (!userExists) {
            console.log('An identity for the user "user1" does not exist in the wallet');
            console.log('Run the registerUser.js application before retrying');
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccpPath, { wallet, identity: 'user1', discovery: { enabled: true, asLocalhost: true } });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork('mychannel');

        // Get the contract from the network.
        const contract = network.getContract('fabcar');

        // Submit the specified transaction.
        // createCar transaction - requires 5 argument, ex: ('createCar', 'CAR12', 'Honda', 'Accord', 'Black', 'Tom')
        // changeCarOwner transaction - requires 2 args , ex: ('changeCarOwner', 'CAR10', 'Dave')
        await contract.submitTransaction('createDeal', 'DEAL12', 'H', 'A', amount, '2019-5-26');
        console.log('Transaction has been submitted');

        // Disconnect from the gateway.
        await gateway.disconnect();

    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        process.exit(1);
    }
}

