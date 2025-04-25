import 'package:mysql_client/mysql_client.dart';

// FirebaseFirestore.instance.collection('players').add(
//   players,
// );

Future<List<Map<String, dynamic>>> retrieveResults() async {
  List<Map<String, dynamic>> data = [];

  final conn = await MySQLConnection.createConnection(
      host: '71.135.198.135',
      port: 3306,
      userName: 'root',
      password: 'Babycakes15!',
      databaseName: 'saute');

  // var conn = await conn.connect(settings);
  await conn.connect();

  var results = await conn.execute('select * from recipes');

  for (var result in results.rows) {
    var image = result.colAt(0);
    var name = result.colAt(1);
    var prepTime = result.colAt(2);
    var cookTime = result.colAt(3);
    var totalTime = result.colAt(4);
    var calories = result.colAt(5);
    var servings = result.colAt(6);
    var ingredients = result.colAt(7);
    var instructions = result.colAt(8);

    Map<String, dynamic> mappedData = {
      "image": image,
      "name": name,
      "prep time": prepTime,
      "cook time": cookTime,
      "total time": totalTime,
      "calories": calories,
      "servings": servings,
      "ingredients": ingredients,
      "instructions": instructions,
    };

    data.add(mappedData);
  }

  return data;
}

Future<void> retrieveResultsFirebase() async {

}

Future<List<String>> retrieveIngredients() async {
  List<String> recipeIngredients = [];

  final conn = await MySQLConnection.createConnection(
      host: '71.135.198.135',
      port: 3306,
      userName: 'root',
      password: 'Babycakes15!',
      databaseName: 'saute');

  // var conn = await conn.connect(settings);
  await conn.connect();

  var results = await conn.execute('select ingredients from recipes');

  for (var result in results.rows) {
    print(result.assoc());

    var ingredients = result.colAt(0);

    var splitIngredients = ingredients!.split(",");

    for (var i = 0; i < splitIngredients.length; i++) {
      recipeIngredients.add(
        splitIngredients[i].toString(),
      );
    }
  }

  return recipeIngredients;
}

Future<Map<String, dynamic>> retrieveRecipe(String recipeName) async {
  Map<String, dynamic> mappedData = {};

  var conn = await MySQLConnection.createConnection(
    host: '71.135.198.135',
    port: 3306,
    userName: 'root',
    password: 'Babycakes15!',
    databaseName: 'saute',
  );

  await conn.connect();

  var results = await conn.execute(
    "SELECT * FROM recipes WHERE name = :name",
    {"name": recipeName},
  );

  for (var result in results.rows) {
    var image = result.colAt(1);
    var name = result.colAt(2);
    var prepTime = result.colAt(3);
    var cookTime = result.colAt(4);
    var totalTime = result.colAt(5);
    var calories = result.colAt(6);
    var servings = result.colAt(7);
    var ingredients = result.colAt(8);
    var instructions = result.colAt(9);

    mappedData = {
      "image": image,
      "name": name,
      "prep time": prepTime,
      "cook time": cookTime,
      "total time": totalTime,
      "calories": calories,
      "servings": servings,
      "ingredients": ingredients,
      "instructions": instructions,
    };
  }

  return mappedData;
}
Future<List<String>> getShoppingList() async {
  List<String> splitList = [];

  var conn = await MySQLConnection.createConnection(
    host: '71.135.198.135',
    port: 3306,
    userName: 'root',
    password: 'Babycakes15!',
    databaseName: 'saute',
  );

  await conn.connect();

  var results = await conn.execute("SELECT * FROM shopping_list");

  for (var result in results.rows) {
    var assoc = result.assoc();
    print(assoc);
    var list = result.colByName("list");
    print(list);
    var listSplit = list!.split('\n');
    for (var item in listSplit) {
      splitList.add(item);
    }
  }

  return splitList;
}

Future<void> updateShoppingList(List<String> shoppingList) async {
  var conn = await MySQLConnection.createConnection(
    host: '71.135.198.135',
    port: 3306,
    userName: 'root',
    password: 'Babycakes15!',
    databaseName: 'saute',
  );

  await conn.connect();

  // Delete existing shopping list
  await conn.execute("DELETE FROM shopping_list");

  createShoppingListTables();

// Insert new shopping list items
  for (var item in shoppingList) {
    await conn.execute("INSERT INTO shopping_list (list) VALUES (?)", [item] as Map<String, dynamic>?);
  }

  await conn.close();
}

Future<List<String>> retrieveImages() async {
  List<String> images = [];

  var conn = await MySQLConnection.createConnection(
      host: '71.135.198.135',
      port: 3306,
      userName: 'root',
      password: 'Babycakes15!',
      databaseName: 'saute'
  );

  await conn.connect();

  var results = await conn.execute('select * from recipes');

  for (var result in results.rows) {
    print(result.assoc());

    var title = result.colByName("image");

    images.add(title.toString());
  }

  return images;
}

Future<List<String>> retrieveTitles() async {
  List<String> titles = [];

  var conn = await MySQLConnection.createConnection(
      host: '71.135.198.135',
      port: 3306,
      userName: 'root',
      password: 'Babycakes15!',
      databaseName: 'saute');

  await conn.connect();

  var results = await conn.execute('select * from recipes');

  for (var result in results.rows) {
    print(result.assoc());

    var title = result.colByName("name");

    titles.add(title.toString());
  }

  return titles;
}

void createRecipesTables() async {
  var conn = await MySQLConnection.createConnection(
      host: '71.135.198.135',
      port: 3306,
      userName: 'root',
      password: 'Babycakes15!',
      databaseName: 'saute');

  await conn.connect();

  var result = await conn.execute(
      'CREATE TABLE recipes (id int NOT NULL AUTO_INCREMENT PRIMARY KEY, image varchar(255), name varchar(255), preptime varchar(255), cooktime varchar(255), totaltime varchar(255), calories varchar(255), servings varchar(255), ingredients LONGTEXT, instructions LONGTEXT)');

  print(result);
}

void createShoppingListTables() async {
  var conn = await MySQLConnection.createConnection(
      host: '71.135.198.135',
      port: 3306,
      userName: 'root',
      password: 'Babycakes15!',
      databaseName: 'saute');

  await conn.connect();

  var result = await conn.execute(
      'CREATE TABLE shopping_list (id int NOT NULL AUTO_INCREMENT PRIMARY KEY, list LONGTEXT)');

  print(result);
}

Future<void> writeRecipes(Map<String, dynamic> data) async {
  final conn = await MySQLConnection.createConnection(
    host: '71.135.198.135',
    port: 3306,
    userName: 'root',
    password: 'Babycakes15!',  // Delete existing shopping list
    databaseName: 'saute',
  );
  await conn.connect();
  await conn.execute(
    'INSERT INTO recipes (image, name, preptime, cooktime, totaltime, calories, servings, ingredients, instructions) VALUES (:image, :name, :preptime, :cooktime, :totaltime, :calories, :servings, :ingredients, :instructions)',
    data,
  );
  await conn.close();
}

// Future<void> writeShoppingList(List<String> items) async {
//   var conn = await MySQLConnection.createConnection(
//     host: '71.135.198.135',
//     port: 3306,
//     userName: 'root',
//     password: 'Babycakes15!',
//     databaseName: 'saute',
//   );

//   await conn.connect();

//   await conn.execute('DELETE FROM shopping_list');

//   for (var item in items) {
//     await conn.execute(
//       'INSERT INTO shopping_list (item) VALUES (?)',
//       [item] as Map<String, dynamic>?,
//     );
//   }

//   await conn.close();
// }

Future<void> writeShoppingList(Map<String, dynamic> items) async {
  var conn = await MySQLConnection.createConnection(
    host: '71.135.198.135',
    port: 3306,
    userName: 'root',
    password: 'Babycakes15!',
    databaseName: 'saute',
  );
  await conn.connect();

  // Clear the existing shopping list
  await conn.execute('DELETE FROM shopping_list');

  // Prepare the SQL statement for inserting items
  var insertQuery = 'INSERT INTO shopping_list (list) VALUES (:list)';

  for (var item in items['list']) {
    final itemOnList = { "list": item };

    await conn.execute(insertQuery, itemOnList);
  }

  await conn.close();
}

void addToShoppingList(Map<String, dynamic> data) async {
  var conn = await MySQLConnection.createConnection(
      host: '71.135.198.135',
      port: 3306,
      userName: 'root',
      password: 'Babycakes15!',
      databaseName: 'saute');

  await conn.connect();

  var result = await conn.execute(
      'update shopping_list set (list) values (:list)', data);

  print(result);
}

void dropShoppingListTable() async {
  var conn = await MySQLConnection.createConnection(
      host: '71.135.198.135',
      port: 3306,
      userName: 'root',
      password: 'Babycakes15!',
      databaseName: 'saute');

  await conn.connect();

  var result = await conn.execute('drop table shopping_list');

  print(result);
}

// Delete a recipe by name
Future<void> deleteRecipe(String name) async {
  var conn = await MySQLConnection.createConnection(
    host: '71.135.198.135',
    port: 3306,
    userName: 'root',
    password: 'Babycakes15!',
    databaseName: 'saute',
  );

  await conn.connect();

  await conn.execute(
    "DELETE FROM recipes WHERE name = :name",
    {"name": name},
  );

  await conn.close();
}

// Search recipes by name or ingredients
Future<List<Map<String, dynamic>>> searchRecipes(String query) async {
  List<Map<String, dynamic>> data = [];

  var conn = await MySQLConnection.createConnection(
    host: '71.135.198.135',
    port: 3306,
    userName: 'root',
    password: 'Babycakes15!',
    databaseName: 'saute',
  );

  await conn.connect();

  var results = await conn.execute(
    "SELECT * FROM recipes WHERE name LIKE :query OR ingredients LIKE :query",
    {"query": "%$query%"},
  );

  for (var result in results.rows) {
    var image = result.colByName('image');
    var name = result.colByName('name');
    var prepTime = result.colByName('preptime');
    var cookTime = result.colByName('cooktime');
    var totalTime = result.colByName('totaltime');
    var calories = result.colByName('calories');
    var servings = result.colByName('servings');
    var ingredients = result.colByName('ingredients');
    var instructions = result.colByName('instructions');

    Map<String, dynamic> mappedData = {
      'image': image,
      'name': name,
      'prep time': prepTime,
      'cook time': cookTime,
      'total time': totalTime,
      'calories': calories,
      'servings': servings,
      'ingredients': ingredients,
      'instructions': instructions,
    };

    data.add(mappedData);
  }

  await conn.close();
  return data;
}