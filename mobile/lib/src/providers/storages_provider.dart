import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:shared_preferences/shared_preferences.dart';

final secureStorageProvider = StateProvider((_) => FlutterSecureStorage());
final sharedPreferencesProvider =
    FutureProvider((_) async => await SharedPreferences.getInstance());
