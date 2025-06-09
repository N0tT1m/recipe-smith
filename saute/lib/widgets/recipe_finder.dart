import 'package:saute/widgets/recipe.dart';
import 'package:flutter/material.dart';
import 'package:advanced_search/advanced_search.dart';
import 'package:saute/services/db.dart';

class RecipeFinder extends StatefulWidget {
  const RecipeFinder({Key? key}) : super(key: key);

  @override
  State<RecipeFinder> createState() => _RecipeFinderState();
}

class _RecipeFinderState extends State<RecipeFinder> {
  late final List<Map<String, dynamic>> recipes;
  late final List<String> titles;
  late final List<String> ingredients;
  late Map<String, dynamic> recipe = {};
  late String recipeTitle = "";

  @override
  void initState() {
    // TODO: implement initState
    setState(() {
      getRecipes();
      getTitles();
      getIngredients();
    });
    super.initState();
  }

  void getRecipes() {
    retrieveResults().then(
          (value) => setState(
            () {
          recipes = value;
        },
      ),
    );
  }

  void getTitles() {
    retrieveTitles().then((value) => setState(() {
      titles.addAll(value);
    }));
  }

  void getIngredients() {
    retrieveIngredients().then((value) => setState(() {
      ingredients = value;
    }));
  }

  void compareIngredientsToRecipes() {
    
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
      body: Container(
        alignment: Alignment.center,
        height: MediaQuery.of(context).size.height,
        width: MediaQuery.of(context).size.width,
        child: recipe == null
            ? AdvancedSearch(
          // This is basically an Input Text Field
          searchItems: ingredients,
          maxElementsToDisplay: 5,
          onItemTap: (index, text) {
            // user just found what he needs, now it's your turn to handle that
          },
          onSearchClear: () {
            // may be display the full list? or Nothing? it's your call
          },
          onSubmitted: (value, value2) {
            setState(() {
              recipeTitle = value;
            });

            retrieveRecipe(recipeTitle).then((value) => recipe = value);
          },
          onEditingProgress: (value, value2) {
            // user is trying to, searchItems: [] lookup something, may be you want to help him?
          },
        ) : Recipe(recipeName: recipeTitle),
      ),
    );
  }
}
