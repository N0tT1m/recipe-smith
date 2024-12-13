import 'package:saute/widgets/recipe.dart';
import 'package:flutter/material.dart';
import 'package:saute/services/db.dart';
import 'package:advanced_search/advanced_search.dart';


class RecipeSearch extends StatefulWidget {
  const RecipeSearch({Key? key}) : super(key: key);

  @override
  State<RecipeSearch> createState() => _RecipeSearchState();
}

class _RecipeSearchState extends State<RecipeSearch> {
  List<String> recipes = [];
  List<String> titles = [];
  var recipe;
  var recipeTitle;

  @override
  void initState() {
    // TODO: implement initState
    setState(() {
      recipes = getRecipes();
      titles = getTitles();
    });
    super.initState();
  }

  List<String> getRecipes() {
    retrieveResults().then(
          (value) => setState(
            () {
          var recipes = value;
        },
      ),
    );

    return recipes;
  }

  List<String> getTitles() {
    retrieveTitles().then((value) => setState(() {
      titles.addAll(value);
    }));

    return titles;
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
          searchItems: titles,
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