import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/features/home/navigators/home_navigator.dart';
import 'package:mobile/src/features/home/navigators/profile_navigator.dart';
import 'package:mobile/src/features/petifies/screens/petifies_home_screen.dart';
import 'package:mobile/src/providers/user_model_providers.dart';
import 'package:mobile/src/widgets/bottom_nav_bars/main_bottom_nav_bar.dart';

enum HomeTabItem { home, search, petifies, notification, profile }

class HomeScreeen extends ConsumerStatefulWidget {
  const HomeScreeen({super.key});

  @override
  ConsumerState<ConsumerStatefulWidget> createState() => _HomeScreeenState();
}

class _HomeScreeenState extends ConsumerState<HomeScreeen> {
  var _selectedTab = HomeTabItem.home;

  final _navigatorKeys = {
    HomeTabItem.home: GlobalKey<NavigatorState>(),
    HomeTabItem.search: GlobalKey<NavigatorState>(),
    HomeTabItem.petifies: GlobalKey<NavigatorState>(),
    HomeTabItem.notification: GlobalKey<NavigatorState>(),
    HomeTabItem.profile: GlobalKey<NavigatorState>(),
  };

  @override
  void initState() {
    super.initState();
  }

  void _onItemTapped(int index) {
    setState(() {
      _selectedTab = HomeTabItem.values[index];
    });
  }

  Widget _buildOffStageHomeNavigator() {
    return Offstage(
      offstage: _selectedTab != HomeTabItem.home,
      child: HomeNavigator(
        navigatorKey: _navigatorKeys[HomeTabItem.home]!,
      ),
    );
  }

  Widget _buildOffStageSearchNavigator() {
    return Offstage(
      offstage: _selectedTab != HomeTabItem.search,
      child: Placeholder(),
    );
  }

  Widget _buildOffStagePetifiesNavigator() {
    return Offstage(
      offstage: _selectedTab != HomeTabItem.petifies,
      child: PetifiesHomeScreeen(),
    );
  }

  Widget _buildOffStageNotificationNavigator() {
    return Offstage(
      offstage: _selectedTab != HomeTabItem.notification,
      child: Placeholder(),
    );
  }

  Widget _buildOffStageProfile() {
    return Offstage(
      offstage: _selectedTab != HomeTabItem.profile,
      child: ProfileNavigator(
        navigatorKey: _navigatorKeys[HomeTabItem.profile]!,
      ),
    );
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
      body: Stack(
        children: [
          _buildOffStageHomeNavigator(),
          _buildOffStageSearchNavigator(),
          _buildOffStagePetifiesNavigator(),
          _buildOffStageNotificationNavigator(),
          _buildOffStageProfile()
        ],
      ),
      bottomNavigationBar: MainButtomNavBar(
        curPage: _selectedTab.index,
        onTapFunc: _onItemTapped,
      ),
    );
  }
}
