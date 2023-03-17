// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:mobile/src/theme/themes.dart';

class MainButtomNavBar extends StatelessWidget {
  final int curPage;
  final Function(int) onTapFunc;

  const MainButtomNavBar({
    Key? key,
    required this.curPage,
    required this.onTapFunc,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        color: Colors.white,
        border: Border(
          top: BorderSide(
            color: Theme.of(context).colorScheme.tertiary,
            width: 0.1,
          ),
        ),
      ),
      child: BottomNavigationBar(
        backgroundColor: Theme.of(context).appBarTheme.backgroundColor,
        unselectedItemColor: Themes.whiteColor,
        showUnselectedLabels: true,
        unselectedFontSize: 10,
        selectedFontSize: 10,
        selectedItemColor: Themes.blueColor,
        type: BottomNavigationBarType.fixed,
        currentIndex: curPage,
        onTap: onTapFunc,
        items: [
          BottomNavigationBarItem(
            icon: Icon(Icons.home),
            label: "Home",
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.search),
            label: "Searching",
          ),
          BottomNavigationBarItem(
            icon: ClipOval(
              child: Material(
                color: Themes.blueColor,
                child: InkWell(
                  splashColor: Theme.of(context).colorScheme.inversePrimary,
                  onTap: () {},
                  child: SizedBox(
                    width: 50,
                    height: 50,
                    child: Icon(
                      Icons.pets,
                    ),
                  ),
                ),
              ),
            ),
            label: "Petifies",
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.notifications),
            label: "Notifications",
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.account_circle),
            label: "Profile",
          ),
        ],
      ),
    );
  }
}
