import 'package:flutter/material.dart';
import 'package:saute/services/db.dart';
import 'package:vertical_card_pager/vertical_card_pager.dart';
import 'package:saute/widgets/recipe.dart';

class Recipes extends StatefulWidget {
  const Recipes({Key? key}) : super(key: key);

  @override
  State<Recipes> createState() => _RecipesState();
}

class _RecipesState extends State<Recipes> {
  List<String> recipes = [];
  List<String> titles = [];
  List<String> images = [];
  List<Widget> widgetImages = [];

  @override
  void initState() {
    // TODO: implement initState
    setState(() {
      recipes = getRecipes();
      titles = getTitles();
      widgetImages = getImages();
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

  List<Widget> getImages() {
    List<Widget> networkImages = [];

    retrieveImages().then((value) => setState(() {
          images.addAll(value);

          for (var image in images) {
            setState(() {
              networkImages.add(Image.network(image));
            });
          }
        }));

    return networkImages;
  }

  @override
  Widget build(BuildContext context) {
    // final List<String> titles = [
    //   "RED",
    //   "YELLOW",
    //   "BLACK",
    //   "CYAN",
    //   "BLUE",
    //   "GREY",
    // ];

    //final List<Widget> images = [
    //  Container(
    //    color: Colors.red,
    // ),
    //];

    return Scaffold(
      appBar: AppBar(
        leading: BackButton(
          color: Theme.of(context).colorScheme.onSecondary,
          onPressed: () => Navigator.of(context).pop(),
        ),
      ),
      body: SafeArea(
        child: Column(
          children: <Widget>[
            Expanded(
              child: Container(
                child: VerticalCardPager(
                  textStyle: const TextStyle(
                      color: Colors.white, fontWeight: FontWeight.bold),
                  titles: titles,
                  images: widgetImages,
                  onPageChanged: (page) {},
                  align: ALIGN.CENTER,
                  onSelectedItem: (index) {
                    Navigator.push(
                      context,
                      MaterialPageRoute(
                        builder: (context) => Recipe(
                          recipeName: titles[index],
                        ),
                      ),
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
