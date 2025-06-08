import 'dart:convert';
import 'package:http/http.dart' as http;

class ESService {
  final String apiBaseUrl;

  ESService({required this.apiBaseUrl});

  // Get recipe from Elasticsearch API
  Future<Map<String, dynamic>> getRecipe(String id) async {
    final response = await http.get(
      Uri.parse('$apiBaseUrl/api/recipes/$id'),
      headers: {'Content-Type': 'application/json'},
    );

    if (response.statusCode == 200) {
      return json.decode(response.body);
    } else {
      throw Exception('Failed to load recipe: ${response.statusCode}');
    }
  }

  // Convert Elasticsearch data to app format
  Map<String, dynamic> convertToAppFormat(Map<String, dynamic> esRecipe) {
    // Make sure all required fields exist with at least empty strings
    return {
      'id': esRecipe['id'] ?? '',
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

  // Helper method to format ingredients/instructions that might be lists or strings
  String _formatListOrString(dynamic value) {
    if (value == null) {
      return '';
    } else if (value is List) {
      return value.join(';');
    } else {
      return value.toString();
    }
  }

  // Search recipes with query
  Future<List<Map<String, dynamic>>> searchRecipes(String query, {int page = 1, int size = 20}) async {
    final response = await http.get(
      Uri.parse('$apiBaseUrl/api/recipes/search?q=${Uri.encodeComponent(query)}&page=$page&size=$size'),
      headers: {'Content-Type': 'application/json'},
    );

    if (response.statusCode == 200) {
      List<dynamic> results = json.decode(response.body);
      return results.map((recipe) => convertToAppFormat(recipe)).toList();
    } else {
      throw Exception('Failed to search recipes: ${response.statusCode}');
    }
  }

  // Get recipes by category or tag
  Future<List<Map<String, dynamic>>> getRecipesByCategory(String category, {int page = 1, int size = 20}) async {
    final response = await http.get(
      Uri.parse('$apiBaseUrl/api/recipes/category/${Uri.encodeComponent(category)}?page=$page&size=$size'),
      headers: {'Content-Type': 'application/json'},
    );

    if (response.statusCode == 200) {
      List<dynamic> results = json.decode(response.body);
      return results.map((recipe) => convertToAppFormat(recipe)).toList();
    } else {
      throw Exception('Failed to get recipes by category: ${response.statusCode}');
    }
  }

  // Get recent recipes - can be used to browse all recipes with pagination
  Future<List<Map<String, dynamic>>> getRecentRecipes({int page = 1, int size = 20}) async {
    final response = await http.get(
      Uri.parse('$apiBaseUrl/api/recipes/recent?page=$page&size=$size'),
      headers: {'Content-Type': 'application/json'},
    );

    if (response.statusCode == 200) {
      List<dynamic> results = json.decode(response.body);
      return results.map((recipe) => convertToAppFormat(recipe)).toList();
    } else {
      throw Exception('Failed to get recent recipes: ${response.statusCode}');
    }
  }

  // Get all recipes - using empty string search to get all recipes
  // This is a workaround since there's no dedicated endpoint for all recipes
  Future<List<Map<String, dynamic>>> getAllRecipes({int page = 1, int size = 50}) async {
    // Using match_all query through an empty string search
    // This is more efficient than using the recent recipes endpoint
    // which has a date filter
    final response = await http.get(
      Uri.parse('$apiBaseUrl/api/recipes/search?q=&page=$page&size=$size'),
      headers: {'Content-Type': 'application/json'},
    );

    if (response.statusCode == 200) {
      List<dynamic> results = json.decode(response.body);
      return results.map((recipe) => convertToAppFormat(recipe)).toList();
    } else {
      throw Exception('Failed to get all recipes: ${response.statusCode}');
    }
  }

  // Parse and format ingredients from semicolon-separated string
  List<String> parseIngredients(String ingredientsString) {
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

  // Parse and format instructions from semicolon-separated string
  List<String> parseInstructions(String instructionsString) {
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