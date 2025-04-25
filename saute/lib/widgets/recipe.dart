import 'package:saute/services/db.dart';
import 'package:flutter/material.dart';

class Recipe extends StatefulWidget {
  final String recipeName;
  const Recipe({Key? key, required this.recipeName}) : super(key: key);

  @override
  State<Recipe> createState() => _RecipeState();
}

class _RecipeState extends State<Recipe> {
  late final ingredients;
  late final instructions;
  Map<String, dynamic> recipe = {};

  @override
  void initState() {
    super.initState();
    getRecipe();
  }

  void getRecipe() {
    retrieveRecipe(widget.recipeName).then((value) {
      setState(() {
        recipe = value;
        ingredients = recipe['ingredients'].split(";").toList();
        instructions = recipe['instructions'].split(";").toList();
      });
    });
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
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(5),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: <Widget>[
            if (recipe.containsKey("image"))
              Image.network(recipe["image"]),
            const Padding(padding: EdgeInsets.all(6)),
            if (recipe.containsKey("name"))
              Text("Recipe Name: ${recipe["name"]}"),
            const Padding(padding: EdgeInsets.all(6)),
            if (recipe.containsKey("prep time"))
              Text("Prep Time: ${recipe["prep time"]}"),
            const Padding(padding: EdgeInsets.all(6)),
            if (recipe.containsKey("cook time"))
              Text("Cook Time: ${recipe["cook time"]}"),
            const Padding(padding: EdgeInsets.all(6)),
            if (recipe.containsKey("total time"))
              Text("Total Time: ${recipe["total time"]}"),
            const Padding(padding: EdgeInsets.all(6)),
            if (recipe.containsKey("calories"))
              Text("Calories: ${recipe["calories"]}"),
            const Padding(padding: EdgeInsets.all(6)),
            if (recipe.containsKey("servings"))
              Text("Servings: ${recipe["servings"]}"),
            const Padding(padding: EdgeInsets.all(6)),
            const Text("Ingredients:"),
            SizedBox(
              height: 200,
              child: Scrollbar(
                child: ListView.builder(
                  shrinkWrap: true,
                  itemCount: ingredients.length - 1,
                  itemBuilder: (context, index) {
                    return ListTile(
                      title: Text('- ${ingredients[index]}'),
                    );
                  },
                ),
              ),
            ),
            const Padding(padding: EdgeInsets.all(6)),
            const Text("Instructions:"),
            SizedBox(
              height: 200,
              child: Scrollbar(
                child: ListView.builder(
                  shrinkWrap: true,
                  itemCount: instructions.length - 1,
                  itemBuilder: (context, index) {
                    return ListTile(
                      title: Text('${index + 1}. ${instructions[index]}'),
                    );
                  },
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}