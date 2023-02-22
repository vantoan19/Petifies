import 'dart:ui';

import 'package:flutter/material.dart';

class Themes {
  // Colors
  static const blackColor = Color.fromRGBO(1, 1, 1, 1); // primary color
  static const greyColor = Color.fromRGBO(90, 90, 90, 1); // secondary color
  static const lightGreyColor = Color.fromRGBO(175, 175, 175, 1);
  static const drawerColor = Color.fromRGBO(18, 18, 18, 1);
  static const whiteColor = Colors.white;
  static const blueColor = Color.fromRGBO(60, 185, 255, 1);
  static const redColor = Colors.red;
  static const yellowColor = Colors.yellow;
  static const splashBlueColor = Color.fromRGBO(118, 206, 255, 1);

  static final textTheme = const TextTheme().apply(fontFamily: 'Roboto');

  // Themes
  static var darkModeAppTheme = ThemeData.dark().copyWith(
    scaffoldBackgroundColor: blackColor,
    cardColor: whiteColor,
    appBarTheme: const AppBarTheme(
      backgroundColor: blackColor,
      iconTheme: IconThemeData(
        color: whiteColor,
      ),
    ),
    drawerTheme: const DrawerThemeData(
      backgroundColor: drawerColor,
    ),
    colorScheme: const ColorScheme.dark().copyWith(
        primary: blueColor,
        inversePrimary: whiteColor,
        secondary: lightGreyColor),
    primaryColor: blueColor,
    primaryTextTheme: textTheme.apply(displayColor: whiteColor),
  );

  static var lightModeAppTheme = ThemeData.light().copyWith(
    scaffoldBackgroundColor: whiteColor,
    cardColor: greyColor,
    appBarTheme: const AppBarTheme(
      backgroundColor: whiteColor,
      elevation: 0,
      iconTheme: IconThemeData(
        color: blackColor,
      ),
    ),
    drawerTheme: const DrawerThemeData(
      backgroundColor: whiteColor,
    ),
    colorScheme: const ColorScheme.light().copyWith(
        primary: blueColor, inversePrimary: blackColor, secondary: greyColor),
    primaryColor: blueColor,
  );
}
