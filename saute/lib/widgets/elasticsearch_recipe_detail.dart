import 'dart:io';

import 'package:flutter/material.dart';
import 'package:saute/services/es_service.dart';
import 'package:saute/services/db.dart';

class ElasticsearchRecipeDetail extends StatefulWidget {
  final String recipeId;

  const ElasticsearchRecipeDetail({
    Key? key,
    required this.recipeId,
  }) : super(key: key);

  @override
  _ElasticsearchRecipeDetailState createState() => _ElasticsearchRecipeDetailState();
}

class _ElasticsearchRecipeDetailState extends State<ElasticsearchRecipeDetail> {
  late Future<Map<String, dynamic>> _recipeFuture;
  bool _isSaving = false;

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

  // Initialize the Elasticsearch service
  final ESService _esService = ESService(
      apiBaseUrl: baseUrl // 'http://10.0.2.2:8080', // Update to your API server address
  );

  @override
  void initState() {
    super.initState();
    _recipeFuture = _loadRecipe();
  }

  // Load recipe details from Elasticsearch
  Future<Map<String, dynamic>> _loadRecipe() async {
    try {
      // Get recipe from Elasticsearch
      final recipe = await _esService.getRecipe(widget.recipeId);

      // Convert to app format
      return _esService.convertToAppFormat(recipe);
    } catch (e) {
      print('Error loading recipe: $e');
      throw e;
    }
  }

  // Save recipe to app database
  Future<void> _saveRecipe(Map<String, dynamic> recipe) async {
    setState(() {
      _isSaving = true;
    });

    try {
      // Format the recipe data for your MySQL database
      // The keys in the data map must match the column names in the database
      final recipeData = {
        'image': recipe['image'],
        'name': recipe['name'],
        'prep_time': recipe['prep_time'],
        'cook_time': recipe['cook_time'],
        'total_time': recipe['total_time'],
        'calories': recipe['calories'],
        'servings': recipe['servings'],
        'ingredients': recipe['ingredients'],
        'instructions': recipe['instructions'],
      };

      // Save to MySQL database
      await writeRecipes(recipeData);

      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('Recipe saved to your collection!'),
          duration: Duration(seconds: 2),
        ),
      );
    } catch (e) {
      print('Error saving recipe: $e');
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('Error saving recipe: $e'),
          duration: Duration(seconds: 2),
        ),
      );
    } finally {
      setState(() {
        _isSaving = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Recipe Details'),
        actions: [
          // Save button in app bar
          FutureBuilder<Map<String, dynamic>>(
            future: _recipeFuture,
            builder: (context, snapshot) {
              if (snapshot.connectionState == ConnectionState.done && snapshot.hasData) {
                return IconButton(
                  icon: _isSaving
                      ? SizedBox(
                    width: 20,
                    height: 20,
                    child: CircularProgressIndicator(
                      color: Colors.white,
                      strokeWidth: 2,
                    ),
                  )
                      : Icon(Icons.save),
                  onPressed: _isSaving
                      ? null
                      : () => _saveRecipe(snapshot.data!),
                  tooltip: 'Save to My Recipes',
                );
              }
              return SizedBox.shrink();
            },
          ),
        ],
      ),
      body: FutureBuilder<Map<String, dynamic>>(
        future: _recipeFuture,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return Center(child: CircularProgressIndicator());
          }

          if (snapshot.hasError) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(Icons.error_outline, size: 48, color: Colors.red),
                  SizedBox(height: 16),
                  Text(
                    'Error loading recipe',
                    style: TextStyle(fontSize: 18),
                  ),
                  SizedBox(height: 8),
                  Text(
                    '${snapshot.error}',
                    style: TextStyle(fontSize: 14, color: Colors.grey[700]),
                  ),
                ],
              ),
            );
          }

          if (!snapshot.hasData || snapshot.data!.isEmpty) {
            return Center(child: Text('Recipe not found'));
          }

          final recipe = snapshot.data!;

          return SingleChildScrollView(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // Recipe image
                recipe['image'] != null && recipe['image'].toString().isNotEmpty
                    ? Image.network(
                  recipe['image'].toString(),
                  width: double.infinity,
                  height: 200,
                  fit: BoxFit.cover,
                  errorBuilder: (context, error, stackTrace) => Container(
                    width: double.infinity,
                    height: 200,
                    color: Colors.grey[300],
                    child: Icon(Icons.restaurant, size: 80),
                  ),
                )
                    : Container(
                  width: double.infinity,
                  height: 200,
                  color: Colors.grey[300],
                  child: Icon(Icons.restaurant, size: 80),
                ),

                Padding(
                  padding: const EdgeInsets.all(16.0),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      // Recipe title
                      Text(
                        recipe['name'].toString(),
                        style: TextStyle(
                          fontSize: 24,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      SizedBox(height: 16),

                      // Recipe source
                      if (recipe['source_site'] != null && recipe['source_site'].toString().isNotEmpty)
                        Chip(
                          label: Text(recipe['source_site'].toString()),
                          backgroundColor: Theme.of(context).colorScheme.primary.withOpacity(0.1),
                        ),
                      SizedBox(height: 16),

                      // Recipe info row
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                        children: [
                          _buildInfoItem(Icons.timer, 'Prep', recipe['prep_time'].toString()),
                          _buildInfoItem(Icons.local_fire_department, 'Cook', recipe['cook_time'].toString()),
                          _buildInfoItem(Icons.timer_3, 'Total', recipe['total_time'].toString()),
                          _buildInfoItem(Icons.people, 'Serves', recipe['servings'].toString()),
                          _buildInfoItem(Icons.whatshot, 'Calories', recipe['calories'].toString()),
                        ],
                      ),

                      SizedBox(height: 24),

                      // Ingredients
                      Text(
                        'Ingredients',
                        style: TextStyle(
                          fontSize: 20,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      SizedBox(height: 8),
                      _buildIngredientsList(recipe['ingredients'].toString()),

                      SizedBox(height: 24),

                      // Instructions
                      Text(
                        'Instructions',
                        style: TextStyle(
                          fontSize: 20,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      SizedBox(height: 8),
                      _buildInstructionsList(recipe['instructions'].toString()),

                      // Original URL link
                      if (recipe['url'] != null &&
                          recipe['url'].toString().isNotEmpty &&
                          !recipe['url'].toString().startsWith('app://'))
                        Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            SizedBox(height: 24),
                            Text(
                              'Source',
                              style: TextStyle(
                                fontSize: 20,
                                fontWeight: FontWeight.bold,
                              ),
                            ),
                            SizedBox(height: 8),
                            InkWell(
                              onTap: () {
                                // Add url_launcher package to open the URL
                                // For now, just print it
                                print('Open URL: ${recipe['url']}');
                              },
                              child: Text(
                                recipe['url'].toString(),
                                style: TextStyle(
                                  color: Colors.blue,
                                  decoration: TextDecoration.underline,
                                ),
                              ),
                            ),
                          ],
                        ),
                    ],
                  ),
                ),
              ],
            ),
          );
        },
      ),
    );
  }

  Widget _buildInfoItem(IconData icon, String label, String value) {
    return Column(
      children: [
        Icon(icon),
        SizedBox(height: 4),
        Text(label, style: TextStyle(fontWeight: FontWeight.bold)),
        Text(value.isEmpty ? 'N/A' : value),
      ],
    );
  }

  Widget _buildIngredientsList(String ingredients) {
    final List<String> items = ingredients.split(';')
        .where((item) => item.trim().isNotEmpty)
        .toList();

    return ListView.builder(
      shrinkWrap: true,
      physics: NeverScrollableScrollPhysics(),
      itemCount: items.length,
      itemBuilder: (context, index) {
        return Padding(
          padding: const EdgeInsets.symmetric(vertical: 4.0),
          child: Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Icon(Icons.fiber_manual_record, size: 12),
              SizedBox(width: 8),
              Expanded(child: Text(items[index].trim())),
            ],
          ),
        );
      },
    );
  }

  Widget _buildInstructionsList(String instructions) {
    final List<String> steps = instructions.split(';')
        .where((step) => step.trim().isNotEmpty)
        .toList();

    return ListView.builder(
      shrinkWrap: true,
      physics: NeverScrollableScrollPhysics(),
      itemCount: steps.length,
      itemBuilder: (context, index) {
        return Padding(
          padding: const EdgeInsets.symmetric(vertical: 8.0),
          child: Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Container(
                width: 24,
                height: 24,
                decoration: BoxDecoration(
                  color: Theme.of(context).primaryColor,
                  shape: BoxShape.circle,
                ),
                child: Center(
                  child: Text(
                    '${index + 1}',
                    style: TextStyle(
                      color: Colors.white,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ),
              ),
              SizedBox(width: 12),
              Expanded(child: Text(steps[index].trim())),
            ],
          ),
        );
      },
    );
  }
}