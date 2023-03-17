import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/features/feed/screens/feed_screen.dart';
import 'package:mobile/src/features/profile/screens/profile_screen.dart';
import 'package:mobile/src/providers/model_providers.dart';
import 'package:mobile/src/widgets/appbars/main_appbar.dart';
import 'package:mobile/src/widgets/bottom_nav_bars/main_bottom_nav_bar.dart';
import 'package:mobile/src/widgets/floating_buttons/new_post_floating_button.dart';

class HomeScreeen extends ConsumerStatefulWidget {
  const HomeScreeen({super.key});

  @override
  ConsumerState<ConsumerStatefulWidget> createState() => _HomeScreeenState();
}

class _HomeScreeenState extends ConsumerState<HomeScreeen> {
  int _selectedIndex = 0;
  PreferredSizeWidget? appBar = MainAppBar();
  Widget body = FeedScreen();
  static Map<String, Map<String, dynamic>> screens = {};

  @override
  void initState() {
    super.initState();
    screens["index_0"] = {"appBar": MainAppBar(), "body": FeedScreen()};
    screens["index_1"] = {"appBar": MainAppBar(), "body": FeedScreen()};
    screens["index_2"] = {"appBar": MainAppBar(), "body": FeedScreen()};
    screens["index_3"] = {"appBar": MainAppBar(), "body": FeedScreen()};
    screens["index_4"] = {
      "appBar": null,
      "body": ProfileScreen(
        navigateCallback: _setScreen,
      )
    };
  }

  void _onItemTapped(int index) {
    setState(() {
      _selectedIndex = index;
      appBar = screens['index_${_selectedIndex}']?['appBar'];
      body = screens['index_${_selectedIndex}']?['body'];
    });
  }

  void _setScreen(PreferredSizeWidget? appBar, Widget body) {
    setState(() {
      this.appBar = appBar;
      this.body = body;
    });
  }

  @override
  Widget build(BuildContext context) {
    final user = ref.watch(myUserProvider);

    final err = user.maybeWhen(
      error: (e, stackTrace) => e.toString(),
      orElse: () => "no err",
    );

    final userInfo = user.when(
      data: (data) => data,
      error: (e, stackTrace) => null,
      loading: () => null,
    );

    return Scaffold(
      appBar: appBar,
      body: body,
      floatingActionButton: NewPostFloatingButton(),
      bottomNavigationBar: MainButtomNavBar(
        curPage: _selectedIndex,
        onTapFunc: _onItemTapped,
      ),
    );
  }
}
