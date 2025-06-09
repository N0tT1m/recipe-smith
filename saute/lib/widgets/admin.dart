import 'package:flutter/material.dart';
import 'package:saute/services/db.dart';

class AdminPanel extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Admin Panel',
      home: AdminPanelPage(),
      theme: ThemeData(
        colorScheme: const ColorScheme(
          brightness: Brightness.dark,
          background: Color.fromRGBO(104, 107, 108, 100),
          error: Color.fromRGBO(224, 32, 31, 100),
          onBackground: Color.fromRGBO(104, 107, 108, 100),
          onError: Color.fromRGBO(224, 32, 31, 100),
          onPrimary: Color.fromRGBO(65, 190, 164, 100),
          onSecondary: Color.fromRGBO(64, 0, 128, 100),
          onSurface: Color.fromRGBO(31, 223, 224, 100),
          primary: Color.fromRGBO(65, 190, 164, 100),
          secondary: Color.fromRGBO(31, 223, 224, 100),
          surface: Color.fromRGBO(31, 223, 224, 100),
        ),
      ),
    );
  }
}

class AdminPanelPage extends StatefulWidget {
  @override
  _AdminPanelPageState createState() => _AdminPanelPageState();
}

class _AdminPanelPageState extends State<AdminPanelPage> {
  final _recipeFormKey = GlobalKey<FormState>();
  final _shoppingListFormKey = GlobalKey<FormState>();
  final _imageController = TextEditingController();
  final _nameController = TextEditingController();
  final _prepTimeController = TextEditingController();
  final _cookTimeController = TextEditingController();
  final _totalTimeController = TextEditingController();
  final _caloriesController = TextEditingController();
  final _servingsController = TextEditingController();
  final List<TextEditingController> _ingredientControllers = [TextEditingController()];
  final List<TextEditingController> _instructionControllers = [TextEditingController()];
  final List<TextEditingController> _shoppingListControllers = [TextEditingController()];

  @override
  void dispose() {
    _imageController.dispose();
    _nameController.dispose();
    _prepTimeController.dispose();
    _cookTimeController.dispose();
    _totalTimeController.dispose();
    _caloriesController.dispose();
    _servingsController.dispose();
    for (var controller in _ingredientControllers) {
      controller.dispose();
    }
    for (var controller in _instructionControllers) {
      controller.dispose();
    }
    for (var controller in _shoppingListControllers) {
      controller.dispose();
    }
    super.dispose();
  }

  void _createTables() {
    createRecipesTables();
    createShoppingListTables();
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text('Success'),
        content: Text('Tables created successfully.'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: Text('OK'),
          ),
        ],
      ),
    );
  }

  void _writeRecipe() {
    if (_recipeFormKey.currentState!.validate()) {
      String image = _imageController.text;
      String name = _nameController.text;
      String prepTime = _prepTimeController.text;
      String cookTime = _cookTimeController.text;
      String totalTime = _totalTimeController.text;
      String calories = _caloriesController.text;
      String servings = _servingsController.text;
      final ingredients = _ingredientControllers.map((controller) => '${controller.text};').join('');
      final instructions = _instructionControllers.map((controller) => '${controller.text};').join('');

      Map<String, dynamic> data = {
        "image": image,
        "name": name,
        "preptime": prepTime,
        "cooktime": cookTime,
        "totaltime": totalTime,
        "calories": calories,
        "servings": servings,
        "ingredients": ingredients,
        "instructions": instructions,
      };

      writeRecipes(data);
      _imageController.clear();
      _nameController.clear();
      _prepTimeController.clear();
      _cookTimeController.clear();
      _totalTimeController.clear();
      _caloriesController.clear();
      _servingsController.clear();
      for (var controller in _ingredientControllers) {
        controller.clear();
      }
      for (var controller in _instructionControllers) {
        controller.clear();
      }

      showDialog(
        context: context,
        builder: (context) => AlertDialog(
          title: Text('Success'),
          content: Text('Recipe added successfully.'),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context),
              child: Text('OK'),
            ),
          ],
        ),
      );
    }
  }

  void _addToShoppingList() {
    if (_shoppingListFormKey.currentState!.validate()) {
      List<String> shoppingList = _shoppingListControllers.map((controller) => controller.text).toList();
      Map<String, dynamic> data = {'list': shoppingList};

      writeShoppingList(data);
      for (var controller in _shoppingListControllers) {
        controller.clear();
      }
      showDialog(
        context: context,
        builder: (context) => AlertDialog(
          title: Text('Success'),
          content: Text('Shopping list updated successfully.'),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context),
              child: Text('OK'),
            ),
          ],
        ),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Theme.of(context).colorScheme.background,
      appBar: AppBar(
        title: Text(
          'Admin Panel',
          style: TextStyle(color: Theme.of(context).colorScheme.primary),
        ),
        backgroundColor: Theme.of(context).colorScheme.background,
        foregroundColor: Theme.of(context).colorScheme.primary,
      ),
      body: SingleChildScrollView(
        padding: EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            ElevatedButton(
              onPressed: _createTables,
              child: Text(
                'Create Tables',
                style: TextStyle(color: Theme.of(context).colorScheme.onSecondary),
              ),
            ),
            SizedBox(height: 16.0),
            Form(
              key: _recipeFormKey,
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  TextFormField(
                    controller: _imageController,
                    decoration: InputDecoration(
                      labelText: 'Image URL',
                      border: OutlineInputBorder(),
                    ),
                    validator: (value) {
                      if (value == null || value.isEmpty) {
                        return 'Please enter an image URL';
                      }
                      return null;
                    },
                  ),
                  SizedBox(height: 16.0),
                  TextFormField(
                    controller: _nameController,
                    decoration: InputDecoration(
                      labelText: 'Recipe Name',
                      border: OutlineInputBorder(),
                    ),
                    validator: (value) {
                      if (value == null || value.isEmpty) {
                        return 'Please enter a recipe name';
                      }
                      return null;
                    },
                  ),
                  SizedBox(height: 16.0),
                  TextFormField(
                    controller: _prepTimeController,
                    decoration: InputDecoration(
                      labelText: 'Prep Time',
                      border: OutlineInputBorder(),
                    ),
                    validator: (value) {
                      if (value == null || value.isEmpty) {
                        return 'Please enter prep time';
                      }
                      return null;
                    },
                  ),
                  SizedBox(height: 16.0),
                  TextFormField(
                    controller: _cookTimeController,
                    decoration: InputDecoration(
                      labelText: 'Cook Time',
                      border: OutlineInputBorder(),
                    ),
                    validator: (value) {
                      if (value == null || value.isEmpty) {
                        return 'Please enter cook time';
                      }
                      return null;
                    },
                  ),
                  SizedBox(height: 16.0),
                  TextFormField(
                    controller: _totalTimeController,
                    decoration: InputDecoration(
                      labelText: 'Total Time',
                      border: OutlineInputBorder(),
                    ),
                    validator: (value) {
                      if (value == null || value.isEmpty) {
                        return 'Please enter total time';
                      }
                      return null;
                    },
                  ),
                  SizedBox(height: 16.0),
                  TextFormField(
                    controller: _caloriesController,
                    decoration: InputDecoration(
                      labelText: 'Calories',
                      border: OutlineInputBorder(),
                    ),
                    validator: (value) {
                      if (value == null || value.isEmpty) {
                        return 'Please enter calories';
                      }
                      return null;
                    },
                  ),
                  SizedBox(height: 16.0),
                  TextFormField(
                    controller: _servingsController,
                    decoration: InputDecoration(
                      labelText: 'Servings',
                      border: OutlineInputBorder(),
                    ),
                    validator: (value) {
                      if (value == null || value.isEmpty) {
                        return 'Please enter servings';
                      }
                      return null;
                    },
                  ),
                  SizedBox(height: 16.0),
                  Text('Ingredients'),
                  ..._ingredientControllers.map((controller) => TextFormField(
                    controller: controller,
                    decoration: InputDecoration(
                      hintText: 'Enter an ingredient',
                      border: OutlineInputBorder(),
                    ),
                    validator: (value) {
                      if (value == null || value.isEmpty) {
                        return 'Please enter an ingredient';
                      }
                      return null;
                    },
                  )).toList(),
                  SizedBox(height: 8.0),
                  ElevatedButton(
                    onPressed: () {
                      setState(() {
                        _ingredientControllers.add(TextEditingController());
                      });
                    },
                    child: Text('Add Ingredient'),
                  ),
                  SizedBox(height: 16.0),
                  Text('Instructions'),
                  ..._instructionControllers.map((controller) => TextFormField(
                    controller: controller,
                    decoration: InputDecoration(
                      hintText: 'Enter an instruction',
                      border: OutlineInputBorder(),
                    ),
                    validator: (value) {
                      if (value == null || value.isEmpty) {
                        return 'Please enter an instruction';
                      }
                      return null;
                    },
                  )).toList(),
                  SizedBox(height: 8.0),
                  ElevatedButton(
                    onPressed: () {
                      setState(() {
                        _instructionControllers.add(TextEditingController());
                      });
                    },
                    child: Text('Add Instruction'),
                  ),
                  SizedBox(height: 16.0),
                  ElevatedButton(
                    onPressed: _writeRecipe,
                    child: Text('Add Recipe'),
                  ),
                ],
              ),
            ),
            SizedBox(height: 16.0),
            Form(
              key: _shoppingListFormKey,
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  Text('Shopping List'),
                  ..._shoppingListControllers.map((controller) => TextFormField(
                    controller: controller,
                    decoration: InputDecoration(
                      hintText: 'Enter a shopping list item',
                      border: OutlineInputBorder(),
                    ),
                    validator: (value) {
                      if (value == null || value.isEmpty) {
                        return 'Please enter a shopping list item';
                      }
                      return null;
                    },
                  )).toList(),
                  SizedBox(height: 8.0),
                  ElevatedButton(
                    onPressed: () {
                      setState(() {
                        _shoppingListControllers.add(TextEditingController());
                      });
                    },
                    child: Text(
                      'Add Shopping List Item',
                      style: TextStyle(color: Theme.of(context).colorScheme.onSecondary),
                    ),
                  ),
                  SizedBox(height: 16.0),
                  ElevatedButton(
                    onPressed: _addToShoppingList,
                    child: Text(
                      'Update Shopping List',
                      style: TextStyle(color: Theme.of(context).colorScheme.onSecondary),
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}