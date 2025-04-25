// ShoppingList.dart

import 'package:flutter/material.dart';
import 'package:saute/services/db.dart';
import 'package:mysql_client/exception.dart';

class ShoppingList extends StatefulWidget {
  final List<String> ingredients;
  final List<String> missingIngredients;

  ShoppingList({
    Key? key,
    this.ingredients = const [],
    this.missingIngredients = const [],
  }) : super(key: key);

  ShoppingList.overloadedContructor({
    Key? key,
    required this.ingredients,
    required this.missingIngredients,
  }) : super(key: key);

  @override
  State<ShoppingList> createState() => _ShoppingListState();
}

class _ShoppingListState extends State<ShoppingList> {
  final TextEditingController _itemController = TextEditingController();
  List<String> _shoppingList = [];

  @override
  void initState() {
    super.initState();
    fetchShoppingList();
  }

  Future<void> fetchShoppingList() async {
    try {
      List<String> shoppingList = await getShoppingList();
      setState(() {
        _shoppingList = shoppingList;
      });
    } on MySQLException catch (e) {
      print('Error fetching shopping list: $e');
    }
  }

  Future<void> writeShoppingList(List<String> shoppingList) async {
    try {
      await updateShoppingList(shoppingList);
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Shopping list saved successfully')),
      );
    } on MySQLException catch (e) {
      print('Error writing shopping list: $e');
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Failed to save shopping list')),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        leading: BackButton(
          color: Theme
              .of(context)
              .colorScheme
              .onSecondary,
          onPressed: () => Navigator.of(context).pop(),
        ),
        title: const Text('Shopping List'),
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            if (widget.ingredients.isNotEmpty)
              Card(
                child: Padding(
                  padding: const EdgeInsets.all(16.0),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        'Ingredients You Have:',
                        style: Theme
                            .of(context)
                            .textTheme
                            .headlineMedium,
                      ),
                      const SizedBox(height: 8.0),
                      ...widget.ingredients.map((ingredient) =>
                          Text(ingredient)),
                    ],
                  ),
                ),
              ),
            if (widget.ingredients.isNotEmpty) const SizedBox(height: 16.0),
            if (widget.missingIngredients.isNotEmpty)
              Card(
                child: Padding(
                  padding: const EdgeInsets.all(16.0),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        'Ingredients You Need:',
                        style: Theme
                            .of(context)
                            .textTheme
                            .headlineMedium,
                      ),
                      const SizedBox(height: 8.0),
                      ...widget.missingIngredients.map((ingredient) =>
                          Text(ingredient)),
                    ],
                  ),
                ),
              ),
            if (widget.missingIngredients.isNotEmpty) const SizedBox(
                height: 16.0),
            Text(
              'Shopping List:',
              style: Theme
                  .of(context)
                  .textTheme
                  .headlineMedium,
            ),
            const SizedBox(height: 8.0),
            Expanded(
              child: ListView.builder(
                itemCount: _shoppingList.length,
                itemBuilder: (context, index) {
                  return ListTile(
                    title: Text(_shoppingList[index]),
                    trailing: IconButton(
                      icon: const Icon(Icons.delete),
                      onPressed: () {
                        setState(() {
                          _shoppingList.removeAt(index);
                        });
                      },
                    ),
                  );
                },
              ),
            ),
            const SizedBox(height: 16.0),
            Row(
              children: [
                Expanded(
                  child: TextField(
                    controller: _itemController,
                    decoration: const InputDecoration(
                      labelText: 'Add an item',
                    ),
                  ),
                ),
                ElevatedButton(
                  onPressed: () {
                    if (_itemController.text.isNotEmpty) {
                      setState(() {
                        _shoppingList.add(_itemController.text);
                        _itemController.clear();
                      });
                    }
                  },
                  child: const Text('Add'),
                ),
              ],
            ),
            const SizedBox(height: 16.0),
            ElevatedButton(
              onPressed: () {
                writeShoppingList(_shoppingList);
              },
              child: const Text('Save Shopping List'),
            ),
          ],
        ),
      ),
    );
  }
}
