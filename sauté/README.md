# Sauté - Recipe Mobile App

Sauté is a cross-platform mobile application built with Flutter that provides an intuitive interface for browsing, searching, and managing recipes from the Recipe Smith collection. It connects to the Sous API server to access the comprehensive recipe database.

## 📱 Features

- **Recipe Search**: Fast full-text search across thousands of recipes
- **Recipe Details**: Complete recipe information with ingredients and instructions
- **Shopping Lists**: Create and manage shopping lists from recipe ingredients
- **Favorites**: Save favorite recipes for quick access
- **Categories**: Browse recipes by category, cuisine, or cooking time
- **Offline Support**: Cache recipes for offline viewing
- **Recipe Scaling**: Adjust ingredient quantities based on servings
- **Cross-Platform**: Native iOS and Android support

## 🚀 Screenshots

*Note: Add screenshots here when available*

## 🛠️ Installation

### Prerequisites
- Flutter SDK 3.x+
- Dart SDK 3.x+
- Android Studio / Xcode (for device testing)
- Recipe Smith API server (Sous) running

### Setup Development Environment

#### 1. Install Flutter
Follow the [official Flutter installation guide](https://docs.flutter.dev/get-started/install) for your platform.

#### 2. Clone and Setup
```bash
# Clone the repository
git clone https://github.com/your-username/recipe-smith.git
cd recipe-smith/sauté

# Install dependencies
flutter pub get

# Check Flutter setup
flutter doctor
```

#### 3. Configure API Connection
Update the API endpoint in `lib/services/db.dart`:
```dart
class DatabaseService {
  static const String baseUrl = 'http://your-api-server:8080';
  // For Android emulator: 'http://10.0.2.2:8080'
  // For iOS simulator: 'http://localhost:8080'
}
```

### Building the App

#### Development Mode
```bash
# Run on connected device/emulator
flutter run

# Run with hot reload
flutter run --hot
```

#### Production Build
```bash
# Build APK for Android
flutter build apk --release

# Build AAB for Google Play Store
flutter build appbundle --release

# Build for iOS
flutter build ios --release
```

## 📱 Platform Support

### Android
- **Minimum SDK**: API 21 (Android 5.0)
- **Target SDK**: API 34 (Android 14)
- **Architecture**: arm64-v8a, armeabi-v7a, x86_64

### iOS
- **Minimum Version**: iOS 12.0
- **Architecture**: arm64, x86_64 (simulator)
- **Devices**: iPhone, iPad

### Supported Features by Platform
| Feature | Android | iOS |
|---------|---------|-----|
| Recipe Search | ✅ | ✅ |
| Offline Cache | ✅ | ✅ |
| Shopping Lists | ✅ | ✅ |
| Share Recipes | ✅ | ✅ |
| Push Notifications | ✅ | ✅ |

## 🎨 User Interface

### Main Screens

#### Home Screen
- Featured recipes
- Quick search bar
- Category shortcuts
- Recent recipes

#### Search Screen
- Advanced search filters
- Search history
- Quick filters (time, difficulty, cuisine)

#### Recipe Detail Screen
- Full recipe information
- Ingredient list with checkboxes
- Step-by-step instructions
- Recipe scaling controls
- Share and favorite options

#### Shopping List Screen
- Organized ingredient lists
- Check off purchased items
- Multiple list management

#### Favorites Screen
- Saved recipes
- Quick access to frequently used recipes

### Design System
- **Primary Colors**: Warm cooking-inspired palette
- **Typography**: Clean, readable fonts optimized for mobile
- **Icons**: Material Design icons with custom recipe icons
- **Layout**: Responsive design for various screen sizes

## 🔧 Configuration

### API Configuration
```dart
// lib/services/db.dart
class DatabaseService {
  static const String baseUrl = 'http://localhost:8080';
  static const int timeoutSeconds = 30;
  static const int cacheExpiryHours = 24;
}
```

### App Settings
```dart
// lib/config/app_config.dart
class AppConfig {
  static const String appName = 'Sauté';
  static const String version = '1.0.0';
  static const int recipesPerPage = 20;
  static const int maxCachedRecipes = 100;
}
```

## 📡 API Integration

### Recipe Service
```dart
class RecipeService {
  // Search recipes
  Future<List<Recipe>> searchRecipes(String query, {int limit = 20, int offset = 0});
  
  // Get recipe by ID
  Future<Recipe> getRecipe(String id);
  
  // Get all recipes with pagination
  Future<List<Recipe>> getAllRecipes({int limit = 20, int offset = 0});
  
  // Get recipe statistics
  Future<RecipeStats> getStats();
}
```

### Offline Support
- **Local Database**: SQLite for offline recipe storage
- **Image Caching**: Cached network images for offline viewing
- **Sync Strategy**: Background sync when network available

## 🗂️ Project Structure

```
sauté/
├── lib/
│   ├── main.dart              # App entry point
│   ├── services/
│   │   └── db.dart           # API service layer
│   ├── widgets/
│   │   ├── add_recipe.dart   # Add recipe widget
│   │   ├── recipe.dart       # Recipe display widget
│   │   ├── recipe_finder.dart # Search widget
│   │   ├── recipe_search.dart # Advanced search
│   │   ├── recipes.dart      # Recipe list widget
│   │   ├── shopping_list.dart # Shopping list widget
│   │   └── admin.dart        # Admin features
│   ├── models/
│   │   ├── recipe.dart       # Recipe data model
│   │   └── shopping_item.dart # Shopping list item model
│   ├── screens/
│   │   ├── home_screen.dart  # Main home screen
│   │   ├── search_screen.dart # Search interface
│   │   └── detail_screen.dart # Recipe details
│   └── utils/
│       ├── constants.dart    # App constants
│       └── helpers.dart      # Utility functions
├── assets/
│   └── images/
│       ├── mushrooms.jpg     # Sample images
│       └── portobello.webp
├── android/                  # Android-specific files
├── ios/                      # iOS-specific files
├── pubspec.yaml             # Dependencies and assets
└── README.md
```

## 🔌 Dependencies

### Core Dependencies
```yaml
dependencies:
  flutter:
    sdk: flutter
  http: ^1.1.0              # HTTP requests
  sqflite: ^2.3.0           # Local database
  shared_preferences: ^2.2.2 # Local storage
  cached_network_image: ^3.3.0 # Image caching
  provider: ^6.1.1          # State management
```

### Development Dependencies
```yaml
dev_dependencies:
  flutter_test:
    sdk: flutter
  flutter_lints: ^3.0.0     # Code linting
  build_runner: ^2.4.7      # Code generation
```

## 🎯 Key Features Implementation

### Recipe Search
```dart
class RecipeSearchWidget extends StatefulWidget {
  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        SearchBar(
          onChanged: (query) => _searchRecipes(query),
        ),
        RecipeGrid(recipes: _searchResults),
      ],
    );
  }
}
```

### Shopping List Management
```dart
class ShoppingListService {
  Future<void> addRecipeToShoppingList(Recipe recipe) async {
    // Parse ingredients and add to shopping list
    final ingredients = recipe.ingredients.split(';');
    for (final ingredient in ingredients) {
      await _addShoppingItem(ingredient);
    }
  }
}
```

### Recipe Scaling
```dart
class RecipeScaler {
  static Recipe scaleRecipe(Recipe recipe, double multiplier) {
    // Scale ingredient quantities
    final scaledIngredients = _scaleIngredients(recipe.ingredients, multiplier);
    return recipe.copyWith(
      ingredients: scaledIngredients,
      servings: (int.parse(recipe.servings) * multiplier).toString(),
    );
  }
}
```

## 🧪 Testing

### Unit Tests
```bash
# Run all tests
flutter test

# Run specific test file
flutter test test/services/recipe_service_test.dart

# Run with coverage
flutter test --coverage
```

### Integration Tests
```bash
# Run integration tests on device
flutter drive --target=test_driver/app.dart
```

### Widget Tests
```dart
testWidgets('Recipe search displays results', (WidgetTester tester) async {
  await tester.pumpWidget(MyApp());
  await tester.enterText(find.byType(TextField), 'chicken');
  await tester.pump();
  
  expect(find.text('Chicken Recipes'), findsOneWidget);
});
```

## 🔧 Development

### Getting Started
1. Ensure Flutter is installed and configured
2. Start the Sous API server
3. Update API endpoint in app configuration
4. Run `flutter pub get` to install dependencies
5. Run `flutter run` to start development

### State Management
The app uses Provider for state management:
```dart
class RecipeProvider extends ChangeNotifier {
  List<Recipe> _recipes = [];
  bool _loading = false;
  
  List<Recipe> get recipes => _recipes;
  bool get loading => _loading;
  
  Future<void> searchRecipes(String query) async {
    _loading = true;
    notifyListeners();
    
    _recipes = await RecipeService.searchRecipes(query);
    _loading = false;
    notifyListeners();
  }
}
```

### Adding New Features
1. Create new widget files in `lib/widgets/`
2. Add navigation routes in `main.dart`
3. Update state management if needed
4. Add corresponding tests
5. Update documentation

## 🚀 Deployment

### Android Deployment

#### Google Play Store
1. Build signed AAB: `flutter build appbundle --release`
2. Upload to Google Play Console
3. Configure store listing and metadata
4. Submit for review

#### Direct APK Distribution
```bash
flutter build apk --release --split-per-abi
```

### iOS Deployment

#### App Store
1. Build iOS release: `flutter build ios --release`
2. Open `ios/Runner.xcworkspace` in Xcode
3. Archive and upload to App Store Connect
4. Configure app metadata and submit for review

#### TestFlight
Use Xcode to distribute beta builds through TestFlight.

## 🐛 Troubleshooting

### Common Issues

#### API Connection Problems
```dart
// Check network connectivity
bool hasConnection = await Connectivity().checkConnectivity() != ConnectivityResult.none;

// Handle timeout errors
try {
  final response = await http.get(url).timeout(Duration(seconds: 30));
} on TimeoutException {
  // Handle timeout
}
```

#### Performance Issues
- Use `ListView.builder` for large lists
- Implement proper image caching
- Optimize API calls with pagination
- Use `const` constructors where possible

#### Platform-Specific Issues
- **Android**: Check `android/app/src/main/AndroidManifest.xml` for permissions
- **iOS**: Verify `ios/Runner/Info.plist` configuration
- **Network**: Add network security config for HTTP connections

## 📱 Future Enhancements

### Planned Features
- [ ] **Recipe Collections**: Organize recipes into custom collections
- [ ] **Meal Planning**: Weekly meal planning with calendar integration
- [ ] **Nutritional Info**: Display nutritional information per recipe
- [ ] **Recipe Timer**: Built-in cooking timers for each step
- [ ] **Voice Commands**: Voice-controlled recipe reading
- [ ] **Social Features**: Share recipes with friends
- [ ] **Recipe Rating**: Rate and review recipes
- [ ] **Grocery Integration**: Connect with grocery delivery services

### Technical Improvements
- [ ] **Background Sync**: Automatic recipe updates
- [ ] **Push Notifications**: New recipe alerts
- [ ] **Biometric Auth**: Secure favorites with fingerprint/face ID
- [ ] **Watch App**: Companion app for Apple Watch/Wear OS
- [ ] **Tablet Optimization**: Enhanced UI for tablets
- [ ] **Dark Mode**: Dark theme support
- [ ] **Accessibility**: Enhanced accessibility features

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/new-feature`)
3. Follow Flutter development best practices
4. Add tests for new functionality
5. Update documentation
6. Submit a pull request

### Development Guidelines
- Follow [Flutter style guide](https://docs.flutter.dev/development/tools/formatting)
- Write comprehensive widget tests
- Use meaningful commit messages
- Keep widgets focused and reusable
- Document public APIs

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Flutter team for the excellent framework
- Recipe websites for sharing amazing recipes
- Open source community for valuable packages
- Contributors who help improve the app