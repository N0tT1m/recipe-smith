import 'package:flutter/material.dart';
import 'package:saute/services/db.dart';
import 'package:mysql_client/exception.dart';

class AddShoppingList extends StatefulWidget {
  const AddShoppingList({Key? key}) : super(key: key);

  @override
  State<AddShoppingList> createState() => _AddShoppingListState();
}

class _AddShoppingListState extends State<AddShoppingList> {
  final _formKey = GlobalKey<FormState>();
  final List<TextEditingController> _controllers = [];

  @override
  void dispose() {
    for (var controller in _controllers) {
      controller.dispose();
    }
    super.dispose();
  }

  void _addItem() {
    setState(() {
      _controllers.add(TextEditingController());
    });
  }

  void _submitList() {
    if (_formKey.currentState!.validate()) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Processing Data')),
      );

      try {
        createShoppingListTables();
      } on MySQLServerException catch (e) {
        print(e);
      }

      _formKey.currentState!.save();

      List<String> items = _controllers
          .where((controller) => controller.text.isNotEmpty)
          .map((controller) => controller.text)
          .toList();

      print(items);
      // writeShoppingList(items);

      setState(() {
        _controllers.clear();
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        leading: BackButton(
          color: Theme.of(context).colorScheme.onSecondary,
          onPressed: () => Navigator.of(context).pop(),
        ),
      ),
      body: SafeArea(
        child: Container(
          padding: const EdgeInsets.all(5),
          margin: const EdgeInsets.all(5),
          child: Form(
            key: _formKey,
            child: ListView.builder(
              itemCount: _controllers.length + 1,
              itemBuilder: (context, index) {
                if (index == _controllers.length) {
                  return ElevatedButton(
                    onPressed: _addItem,
                    child: const Text('Add Item'),
                  );
                }
                return TextFormField(
                  controller: _controllers[index],
                  validator: (value) {
                    if (value == null || value.isEmpty) {
                      return 'Please enter an item';
                    }
                    return null;
                  },
                  decoration: InputDecoration(
                    labelText: 'Item ${index + 1}',
                  ),
                );
              },
            ),
          ),
        ),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: _submitList,
        child: const Icon(Icons.save),
      ),
    );
  }
}