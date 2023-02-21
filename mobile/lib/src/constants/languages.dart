import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:flutter_gen/gen_l10n/petifies_localizations.dart';

class Language {
  static SharedPreferences? _prefs = null;
  static const String LANGUAGE_CODE = 'languageCode';

  static const String ENGLISH = 'en';
  static const String VIETNAMESE = 'vi';
  static const String HUNGARIAN = 'hu';

  static Locale _getLocale(String languageCode) {
    switch (languageCode) {
      case ENGLISH:
        return Locale(ENGLISH, 'US');
      case VIETNAMESE:
        return Locale(VIETNAMESE, 'VN');
      case HUNGARIAN:
        return Locale(HUNGARIAN, 'HU');
      default:
        return Locale(ENGLISH, 'US');
    }
  }

  static Future<Locale> setLocale(String languageCode) async {
    if (_prefs == null) {
      _prefs = await SharedPreferences.getInstance();
    }
    await _prefs?.setString(LANGUAGE_CODE, languageCode);
    return _getLocale(languageCode);
  }

  static Future<Locale> getLocale() async {
    if (_prefs == null) {
      _prefs = await SharedPreferences.getInstance();
    }
    String languageCode = _prefs?.getString(LANGUAGE_CODE) ?? ENGLISH;
    return _getLocale(languageCode);
  }

  static AppLocalizations translate(BuildContext context) {
    return AppLocalizations.of(context)!;
  }
}
