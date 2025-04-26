import 'dart:io' show Platform;
import 'package:saute/widgets/add_recipe.dart';
import 'package:saute/widgets/recipe_finder.dart';
import 'package:saute/widgets/shopping_list.dart';
import 'package:flutter/material.dart';
import 'package:saute/widgets/recipes.dart';
import 'package:saute/widgets/admin.dart';
import 'package:saute/widgets/recipe_search.dart';
import 'package:saute/widgets/elasticsearch_recipes.dart';
import 'package:saute/services/db.dart';

void main() async {
  // await Firebase.initializeApp(
  //   options: DefaultFirebaseOptions.currentPlatform,
  // );

  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Sauté',
      theme: ThemeData(
        // This is the theme of your application.
        //
        // Try running your application with "flutter run". You'll see the
        // application has a blue toolbar. Then, without quitting the app, try
        // changing the primarySwatch below to Colors.green and then invoke
        // "hot reload" (press "r" in the console where you ran "flutter run",
        // or simply save your changes to "hot reload" in a Flutter IDE).
        // Notice that the counter didn't reset back to zero; the application
        // is not restarted.
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
            surface: Color.fromRGBO(31, 223, 224, 100)),
      ),
      home: const MyHomePage(title: 'Sauté Home Page'),
    );
  }
}

class MyHomePage extends StatefulWidget {
  const MyHomePage({super.key, required this.title});

  // This widget is the home page of your application. It is stateful, meaning
  // that it has a State object (defined below) that contains fields that affect
  // how it looks.

  // This class is the configuration for the state. It holds the values (in this
  // case the title) provided by the parent (in this case the App widget) and
  // used by the build method of the State. Fields in a Widget subclass are
  // always marked "final".

  final String title;

  @override
  State<MyHomePage> createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {
  @override
  void initState() {
    // TODO: implement initState
    setState(() {

    });
    createRecipesTables();
    createShoppingListTables();
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    // This method is rerun every time setState is called, for instance as done
    // by the _incrementCounter method above.
    //
    // The Flutter framework has been optimized to make rerunning build methods
    // fast, so that you can just rebuild anything that needs updating rather
    // than having to individually change instances of widgets.
    return Platform.isWindows || Platform.isMacOS || Platform.isLinux || Platform.isFuchsia ? AdminPanel() : Scaffold(
      backgroundColor: Theme.of(context).colorScheme.background,
      appBar: AppBar(
        backgroundColor: Theme.of(context).colorScheme.background,
        foregroundColor: Theme.of(context).colorScheme.primary,
        elevation: 5,
        title: Text("Sauté",
            style: TextStyle(color: Theme.of(context).colorScheme.primary)),
      ),
      drawer: Drawer(
        child: ListView(
          // Important: Remove any padding from the ListView.
          padding: EdgeInsets.zero,
          children: [
            DrawerHeader(
              decoration: BoxDecoration(
                color: Theme.of(context).colorScheme.secondary,
              ),
              child: Image.asset(
                "assets/images/mushrooms.jpg",
                width: MediaQuery.of(context).size.width,
                height: MediaQuery.of(context).size.height,
              ),
            ),
            ListTile(
              title: Text('Add Recipes', style: TextStyle(color: Theme.of(context).colorScheme.onSecondary,)),
              onTap: () {
                Navigator.push(
                  context,
                  MaterialPageRoute(builder: (context) => AddRecipeForm()),
                );
              },
            ),
            ListTile(
              title: Text('My Recipes', style: TextStyle(color: Theme.of(context).colorScheme.onSecondary,)),
              onTap: () {
                Navigator.push(
                  context,
                  MaterialPageRoute(builder: (context) => const Recipes()),
                );
              },
            ),
            ListTile(
              title: Text('Discover Recipes', style: TextStyle(color: Theme.of(context).colorScheme.onSecondary,)),
              onTap: () {
                Navigator.push(
                  context,
                  MaterialPageRoute(builder: (context) => const RecipeSearchScreen()),
                );
              },
            ),
            // Add this new ListTile for browsing all recipes
            ListTile(
              title: Text('Browse All Recipes', style: TextStyle(color: Theme.of(context).colorScheme.onSecondary,)),
              onTap: () {
                Navigator.push(
                  context,
                  MaterialPageRoute(builder: (context) => const BrowseRecipesScreen()),
                );
              },
            ),
            ListTile(
              title: Text('Shopping List', style: TextStyle(color: Theme.of(context).colorScheme.onSecondary,)),
              onTap: () {
                Navigator.push(
                  context,
                  MaterialPageRoute(builder: (context) => ShoppingList()),
                );
              },
            ),
            ListTile(
              title: Text('Recipe Search', style: TextStyle(color: Theme.of(context).colorScheme.onSecondary,)),
              onTap: () {
                Navigator.push(
                  context,
                  MaterialPageRoute(builder: (context) => RecipeSearch()),
                );
              },
            ),
            ListTile(
              title: Text('Recipe Finder', style: TextStyle(color: Theme.of(context).colorScheme.onSecondary,)),
              onTap: () {
                Navigator.push(
                  context,
                  MaterialPageRoute(builder: (context) => const RecipeFinder()),
                );
              },
            ),
          ],
        ),
      ),
      body: Center(
        // Center is a layout widget. It takes a single child and positions it
        // in the middle of the parent.
        child: Column(
          // Column is also a layout widget. It takes a list of children and
          // arranges them vertically. By default, it sizes itself to fit its
          // children horizontally, and tries to be as tall as its parent.
          //
          // Invoke "debug painting" (press "p" in the console, choose the
          // "Toggle Debug Paint" action from the Flutter Inspector in Android
          // Studio, or the "Toggle Debug Paint" command in Visual Studio Code)
          // to see the wireframe for each widget.
          //
          // Column has various properties to control how it sizes itself and
          // how it positions its children. Here we use mainAxisAlignment to
          // center the children vertically; the main axis here is the vertical
          // axis because Columns are vertical (the cross axis would be
          // horizontal).
          mainAxisAlignment: MainAxisAlignment.center,
          children: <Widget>[
            Image.asset("assets/images/portobello.webp")
          ],
        ),
      ),
    );
  }
}