import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:io';

import 'elasticsearch_recipe_detail.dart';

class ElasticsearchRecipes {
  static String get baseUrl {
    if (Platform.isAndroid) {
      // Use 10.0.2.2 for Android emulator, or your computer's network IP for physical device
      const bool isEmulator = bool.fromEnvironment('IS_EMULATOR', defaultValue: true);
      return isEmulator
          ? 'http://10.0.2.2:8080'  // Android emulator - connects to host machine
          : 'http://192.168.1.78:8080';  // Physical device - use your computer's IP
    } else if (Platform.isIOS) {
      // iOS simulator uses localhost, physical device uses your computer's IP
      const bool isSimulator = bool.fromEnvironment('IS_SIMULATOR', defaultValue: true);
      return isSimulator
          ? 'http://localhost:8080'  // iOS simulator
          : 'http://192.168.1.78:8080';  // Physical device - use your computer's IP
    }
    // Fallback
    return 'http://192.168.1.78:8080';
  }

  // Modify the convertToAppFormat method to include the ID
  static Map<String, dynamic> convertToAppFormat(Map<String, dynamic> esRecipe) {
    // Make sure all required fields exist with at least empty strings
    return {
      'id': esRecipe['id'] ?? '', // Include the ID
      'name': esRecipe['title'] ?? 'Unnamed Recipe',
      'image': esRecipe['image'] ?? '',
      'prep_time': esRecipe['prep_time'] ?? '',
      'cook_time': esRecipe['cook_time'] ?? '',
      'total_time': esRecipe['total_time'] ?? '',
      'calories': esRecipe['calories'] ?? '',
      'servings': esRecipe['servings'] ?? '',
      // Handle ingredients and instructions which might be lists or strings
      'ingredients': _formatListOrString(esRecipe['ingredients']),
      'instructions': _formatListOrString(esRecipe['instructions']),
      'url': esRecipe['url'] ?? '',
      'source_site': esRecipe['source_site'] ?? '',
    };
  }

  // Add this new method to get all recipes
  static Future<List<Map<String, dynamic>>> getAllRecipes({int page = 1, int size = 20}) async {
    try {
      final url = Uri.parse('$baseUrl/api/recipes/all?page=$page&size=$size');

      final response = await http.get(
        url,
        headers: {'Content-Type': 'application/json'},
      );

      if (response.statusCode == 200) {
        List<dynamic> results = json.decode(response.body);
        return results.map((item) => convertToAppFormat(item)).toList();
      } else {
        throw Exception('Failed to get all recipes: ${response.statusCode}');
      }
    } catch (e) {
      print('Error getting all recipes: $e');
      throw e;
    }
  }

  // Search for recipes in Elasticsearch
  static Future<List<Map<String, dynamic>>> searchRecipes(String query, {int page = 1, int size = 10}) async {
    try {
      final url = Uri.parse('$baseUrl/api/recipes/search?q=${Uri.encodeComponent(query)}&page=$page&size=$size');

      final response = await http.get(
        url,
        headers: {'Content-Type': 'application/json'},
      );

      if (response.statusCode == 200) {
        List<dynamic> results = json.decode(response.body);
        return results.map((item) => convertToAppFormat(item)).toList();
      } else {
        throw Exception('Failed to search recipes: ${response.statusCode}');
      }
    } catch (e) {
      print('Error searching recipes: $e');
      throw e;
    }
  }

  // Get a specific recipe by ID
  static Future<Map<String, dynamic>> getRecipe(String id) async {
    try {
      final url = Uri.parse('$baseUrl/api/recipes/$id');

      final response = await http.get(
        url,
        headers: {'Content-Type': 'application/json'},
      );

      if (response.statusCode == 200) {
        Map<String, dynamic> result = json.decode(response.body);
        return convertToAppFormat(result);
      } else {
        throw Exception('Failed to get recipe: ${response.statusCode}');
      }
    } catch (e) {
      print('Error getting recipe: $e');
      throw e;
    }
  }

  // Get recipes by category
  static Future<List<Map<String, dynamic>>> getRecipesByCategory(String category, {int page = 1, int size = 10}) async {
    try {
      final url = Uri.parse('$baseUrl/api/recipes/category/${Uri.encodeComponent(category)}?page=$page&size=$size');

      final response = await http.get(
        url,
        headers: {'Content-Type': 'application/json'},
      );

      if (response.statusCode == 200) {
        List<dynamic> results = json.decode(response.body);
        return results.map((item) => convertToAppFormat(item)).toList();
      } else {
        throw Exception('Failed to get recipes by category: ${response.statusCode}');
      }
    } catch (e) {
      print('Error getting recipes by category: $e');
      throw e;
    }
  }

  // Get recent recipes
  static Future<List<Map<String, dynamic>>> getRecentRecipes({int page = 1, int size = 10}) async {
    try {
      final url = Uri.parse('$baseUrl/api/recipes/recent?page=$page&size=$size');

      final response = await http.get(
        url,
        headers: {'Content-Type': 'application/json'},
      );

      if (response.statusCode == 200) {
        List<dynamic> results = json.decode(response.body);
        return results.map((item) => convertToAppFormat(item)).toList();
      } else {
        throw Exception('Failed to get recent recipes: ${response.statusCode}');
      }
    } catch (e) {
      print('Error getting recent recipes: $e');
      throw e;
    }
  }

  // Helper method to format ingredients/instructions that might be lists or strings
  static String _formatListOrString(dynamic value) {
    if (value == null) {
      return '';
    } else if (value is List) {
      return value.join(';');
    } else {
      return value.toString();
    }
  }

  // Parse semicolon-separated ingredients into a list
  static List<String> parseIngredients(String ingredientsString) {
    if (ingredientsString.isEmpty) {
      return [];
    }

    // Split by semicolons and trim each item
    return ingredientsString
        .split(';')
        .map((item) => item.trim())
        .where((item) => item.isNotEmpty)
        .toList();
  }

  // Parse semicolon-separated instructions into a list
  static List<String> parseInstructions(String instructionsString) {
    if (instructionsString.isEmpty) {
      return [];
    }

    // Split by semicolons and trim each item
    return instructionsString
        .split(';')
        .map((item) => item.trim())
        .where((item) => item.isNotEmpty)
        .toList();
  }
}

// New widget for browsing all recipes
class BrowseRecipesScreen extends StatefulWidget {
  const BrowseRecipesScreen({Key? key}) : super(key: key);

  @override
  _BrowseRecipesScreenState createState() => _BrowseRecipesScreenState();
}

class _BrowseRecipesScreenState extends State<BrowseRecipesScreen> {
  List<Map<String, dynamic>> _recipes = [];
  bool _isLoading = true;
  bool _hasMoreRecipes = true;
  int _currentPage = 1;
  final int _pageSize = 20;
  final ScrollController _scrollController = ScrollController();

  @override
  void initState() {
    super.initState();
    _loadRecipes();
    _scrollController.addListener(_scrollListener);
  }

  void _scrollListener() {
    if (_scrollController.position.pixels >= _scrollController.position.maxScrollExtent * 0.8 &&
        !_isLoading &&
        _hasMoreRecipes) {
      _loadMoreRecipes();
    }
  }

  Future<void> _loadRecipes() async {
    setState(() {
      _isLoading = true;
    });

    try {
      final recipes = await ElasticsearchRecipes.getAllRecipes(page: 1, size: _pageSize);
      setState(() {
        _recipes = recipes;
        _isLoading = false;
        _hasMoreRecipes = recipes.length == _pageSize;
        _currentPage = 1;
      });
    } catch (e) {
      setState(() {
        _isLoading = false;
      });
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Error loading recipes: $e')),
      );
    }
  }

  Future<void> _loadMoreRecipes() async {
    if (_isLoading) return;

    setState(() {
      _isLoading = true;
    });

    try {
      final moreRecipes = await ElasticsearchRecipes.getAllRecipes(
        page: _currentPage + 1,
        size: _pageSize,
      );

      setState(() {
        _recipes.addAll(moreRecipes);
        _isLoading = false;
        _hasMoreRecipes = moreRecipes.length == _pageSize;
        _currentPage++;
      });
    } catch (e) {
      setState(() {
        _isLoading = false;
      });
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Error loading more recipes: $e')),
      );
    }
  }

  Future<void> _refreshRecipes() async {
    await _loadRecipes();
    return;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('All Recipes'),
      ),
      body: _isLoading && _recipes.isEmpty
          ? Center(child: CircularProgressIndicator())
          : RefreshIndicator(
        onRefresh: _refreshRecipes,
        child: Column(
          children: [
            Expanded(
              child: _recipes.isEmpty
                  ? Center(child: Text('No recipes found'))
                  : RecipeGridView(
                recipes: _recipes,
                scrollController: _scrollController,
                onRecipeTap: (id) {
                  Navigator.push(
                    context,
                    MaterialPageRoute(
                      builder: (context) => ElasticsearchRecipeDetail(
                        recipeId: id,
                      ),
                    ),
                  );
                },
              ),
            ),
            if (_isLoading && _recipes.isNotEmpty)
              Padding(
                padding: const EdgeInsets.all(8.0),
                child: Center(child: CircularProgressIndicator()),
              ),
          ],
        ),
      ),
    );
  }

  @override
  void dispose() {
    _scrollController.dispose();
    super.dispose();
  }
}

// Modify RecipeGridView to have a more balanced aspect ratio
class RecipeGridView extends StatelessWidget {
  final List<Map<String, dynamic>> recipes;
  final Function(String) onRecipeTap;
  final ScrollController? scrollController;

  const RecipeGridView({
    Key? key,
    required this.recipes,
    required this.onRecipeTap,
    this.scrollController,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return GridView.builder(
      controller: scrollController,
      padding: const EdgeInsets.all(8.0),
      gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: 2,
        childAspectRatio: 1.1, // Higher ratio makes cards shorter
        crossAxisSpacing: 10,
        mainAxisSpacing: 10,
      ),
      itemCount: recipes.length,
      itemBuilder: (context, index) {
        final recipe = recipes[index];
        return RecipeCard(
          recipe: recipe,
          onTap: () => onRecipeTap(recipe['id']),
        );
      },
    );
  }
}

// Modified RecipeCard to be more compact with minimal height
class RecipeCard extends StatelessWidget {
  final Map<String, dynamic> recipe;
  final VoidCallback onTap;

  const RecipeCard({
    Key? key,
    required this.recipe,
    required this.onTap,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Card(
      clipBehavior: Clip.antiAlias,
      elevation: 5,
      margin: EdgeInsets.zero, // Remove default card margin
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(8),
      ),
      child: InkWell(
        onTap: onTap,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          mainAxisSize: MainAxisSize.min, // Keep column as small as possible
          children: [
            // Recipe image with smaller ratio to reduce height
            AspectRatio(
              aspectRatio: 16 / 10, // Slightly taller to reduce overall card height
              child: recipe['image'] != null && recipe['image'].toString().isNotEmpty
                  ? Image.network(
                recipe['image'],
                fit: BoxFit.cover,
                errorBuilder: (context, error, stackTrace) => Container(
                  color: Colors.grey[300],
                  child: Icon(Icons.restaurant, size: 40),
                ),
              )
                  : Container(
                color: Colors.grey[300],
                child: Icon(Icons.restaurant, size: 40),
              ),
            ),

            // Minimal content area
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 8.0, vertical: 5.0), // Slightly increased vertical padding
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisSize: MainAxisSize.min,
                children: [
                  // Recipe title - single line to save space
                  Text(
                    recipe['name'],
                    style: TextStyle(
                      fontWeight: FontWeight.bold,
                      fontSize: 17, // Smaller font
                      color: Colors.blue[300],
                    ),
                    maxLines: 1, // Only one line
                    overflow: TextOverflow.ellipsis,
                  ),

                  // Optional time and source in single row if available
                  if (recipe['total_time'] != null || recipe['source_site'] != null)
                    Padding(
                      padding: const EdgeInsets.only(top: 2.0),
                      child: Row(
                        children: [
                          if (recipe['total_time'] != null && recipe['total_time'].toString().isNotEmpty)
                            Expanded(
                              child: Text(
                                recipe['total_time'],
                                style: TextStyle(
                                  fontSize: 15, // Increased from 10
                                  color: Colors.grey[500],
                                ),
                                maxLines: 1,
                                overflow: TextOverflow.ellipsis,
                              ),
                            ),
                          if (recipe['source_site'] != null && recipe['source_site'].toString().isNotEmpty)
                            Expanded(
                              child: Text(
                                recipe['source_site'],
                                style: TextStyle(
                                  fontSize: 15, // Increased from 10
                                  color: Colors.grey[600],
                                ),
                                textAlign: TextAlign.end,
                                maxLines: 1,
                                overflow: TextOverflow.ellipsis,
                              ),
                            ),
                        ],
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

// Widget for recipe search screen
class RecipeSearchScreen extends StatefulWidget {
  const RecipeSearchScreen({Key? key}) : super(key: key);

  @override
  _RecipeSearchScreenState createState() => _RecipeSearchScreenState();
}

class _RecipeSearchScreenState extends State<RecipeSearchScreen> {
  final TextEditingController _searchController = TextEditingController();
  List<Map<String, dynamic>> _searchResults = [];
  bool _isLoading = false;
  bool _hasSearched = false;

  void _performSearch(String query) async {
    if (query.trim().isEmpty) return;

    setState(() {
      _isLoading = true;
      _hasSearched = true;
    });

    try {
      final results = await ElasticsearchRecipes.searchRecipes(query);
      setState(() {
        _searchResults = results;
        _isLoading = false;
      });
    } catch (e) {
      setState(() {
        _isLoading = false;
      });

      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Error searching recipes: $e')),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Search Recipes'),
      ),
      body: Column(
        children: [
          Padding(
            padding: const EdgeInsets.all(16.0),
            child: TextField(
              controller: _searchController,
              decoration: InputDecoration(
                hintText: 'Search recipes...',
                prefixIcon: Icon(Icons.search),
                suffixIcon: IconButton(
                  icon: Icon(Icons.clear),
                  onPressed: () {
                    _searchController.clear();
                  },
                ),
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(10),
                ),
              ),
              onSubmitted: _performSearch,
            ),
          ),

          Expanded(
            child: _isLoading
                ? Center(child: CircularProgressIndicator())
                : !_hasSearched
                ? Center(child: Text('Search for recipes to get started'))
                : _searchResults.isEmpty
                ? Center(child: Text('No recipes found'))
                : RecipeGridView(
              recipes: _searchResults,
              onRecipeTap: (id) {
                Navigator.push(
                  context,
                  MaterialPageRoute(
                    builder: (context) => ElasticsearchRecipeDetail(
                      recipeId: id,
                    ),
                  ),
                );
              },
            ),
          ),
        ],
      ),
    );
  }

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }
}