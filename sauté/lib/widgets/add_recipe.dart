import 'package:saute/services/db.dart';
import 'package:flutter/material.dart';

class AddRecipeForm extends StatefulWidget {
  @override
  _AddRecipeFormState createState() => _AddRecipeFormState();
}

class _AddRecipeFormState extends State<AddRecipeForm> {
  final _formKey = GlobalKey<FormState>();
  final List<TextEditingController> _controllers = List.generate(7, (_) => TextEditingController());
  final List<TextEditingController> _ingredientControllers = [TextEditingController()];
  final List<TextEditingController> _instructionControllers = [TextEditingController()];

  void _addIngredientField() {
    setState(() {
      _ingredientControllers.add(TextEditingController());
    });
  }

  void _addInstructionField() {
    setState(() {
      _instructionControllers.add(TextEditingController());
    });
  }

  void _submitForm() async {
    final ingredients = _ingredientControllers.map((controller) => '${controller.text};').join('');
    final instructions = _instructionControllers.map((controller) => '${controller.text};').join('');

    print(ingredients);

    if (_formKey.currentState!.validate()) {

      final data = {
        'image': _controllers[0].text,
        'name': _controllers[1].text,
        'preptime': _controllers[2].text,
        'cooktime': _controllers[3].text,
        'totaltime': _controllers[4].text,
        'calories': _controllers[5].text,
        'servings': _controllers[6].text,
        'ingredients': ingredients,
        'instructions': instructions,
      };
      await writeRecipes(data);
      _formKey.currentState!.reset();
      _ingredientControllers.clear();
      _instructionControllers.clear();
      _ingredientControllers.add(TextEditingController());
      _instructionControllers.add(TextEditingController());
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Add Recipe'),
      ),
      body: Padding(
        padding: EdgeInsets.all(16.0),
        child: Form(
          key: _formKey,
          child: ListView(
            children: [
              Card(
                margin: EdgeInsets.fromLTRB(8, 2, 8, 8),
                child: Padding(
                  padding: EdgeInsets.all(10),
                  child: Text(
                    "NOTE: Each ingredient and instruction will be displayed on a new line inside the recipe.",
                    style: TextStyle(fontSize: 16),
                  ),
                ),
              ),
              SizedBox(height: 16),
              TextFormField(
                controller: _controllers[0],
                decoration: InputDecoration(labelText: 'Image URL'),
                validator: (value) => value!.isEmpty ? 'Please enter an image URL' : null,
              ),
              SizedBox(height: 16),
              TextFormField(
                controller: _controllers[1],
                decoration: InputDecoration(labelText: 'Recipe Name'),
                validator: (value) => value!.isEmpty ? 'Please enter a recipe name' : null,
              ),
              SizedBox(height: 16),
              TextFormField(
                controller: _controllers[2],
                decoration: InputDecoration(labelText: 'Prep Time'),
                validator: (value) => value!.isEmpty ? 'Please enter the prep time' : null,
              ),
              SizedBox(height: 16),
              TextFormField(
                controller: _controllers[3],
                decoration: InputDecoration(labelText: 'Cook Time'),
                validator: (value) => value!.isEmpty ? 'Please enter the cook time' : null,
              ),
              SizedBox(height: 16),
              TextFormField(
                controller: _controllers[4],
                decoration: InputDecoration(labelText: 'Total Time'),
                validator: (value) => value!.isEmpty ? 'Please enter the total time' : null,
              ),
              SizedBox(height: 16),
              TextFormField(
                controller: _controllers[5],
                decoration: InputDecoration(labelText: 'Calories'),
                validator: (value) => value!.isEmpty ? 'Please enter the calories' : null,
              ),
              SizedBox(height: 16),
              TextFormField(
                controller: _controllers[6],
                decoration: InputDecoration(labelText: 'Servings'),
                validator: (value) => value!.isEmpty ? 'Please enter the servings' : null,
              ),
              SizedBox(height: 16),
              Text('Ingredients'),
              ..._ingredientControllers.map((controller) => Padding(
                padding: EdgeInsets.only(bottom: 8),
                child: Row(
                  children: [
                    Expanded(
                      child: TextFormField(
                        controller: controller,
                        decoration: InputDecoration(
                          hintText: 'Enter an ingredient',
                        ),
                        validator: (value) => value!.isEmpty ? 'Please enter an ingredient' : null,
                      ),
                    ),
                    IconButton(
                      icon: Icon(Icons.add),
                      onPressed: _addIngredientField,
                    ),
                  ],
                ),
              )),
              SizedBox(height: 16),
              Text('Instructions'),
              ..._instructionControllers.map((controller) => Padding(
                padding: EdgeInsets.only(bottom: 8),
                child: Row(
                  children: [
                    Expanded(
                      child: TextFormField(
                        controller: controller,
                        decoration: InputDecoration(
                          hintText: 'Enter an instruction',
                        ),
                        validator: (value) => value!.isEmpty ? 'Please enter an instruction' : null,
                      ),
                    ),
                    IconButton(
                      icon: Icon(Icons.add),
                      onPressed: _addInstructionField,
                    ),
                  ],
                ),
              )),
              SizedBox(height: 16),
              ElevatedButton(
                onPressed: _submitForm,
                child: Text('Submit'),
              ),
            ],
          ),
        ),
      ),
    );
  }
}