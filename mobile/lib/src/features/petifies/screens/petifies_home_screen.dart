import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/features/home/navigators/profile_navigator.dart';
import 'package:mobile/src/features/petifies/navigators/explore_navigator.dart';
import 'package:mobile/src/features/petifies/navigators/my_petifies_navigator.dart';
import 'package:mobile/src/features/petifies/screens/petifies_dashboard.dart';
import 'package:mobile/src/features/petifies/screens/petifies_explore_screen.dart';
import 'package:mobile/src/providers/user_model_providers.dart';
import 'package:mobile/src/widgets/appbars/petifies_main_appbar.dart';
import 'package:mobile/src/widgets/bottom_nav_bars/petifies_bottom_nav_bar.dart';

enum PetifiesHomeTabItem {
  explore,
  mypetifies,
  newpetify,
  notification,
  profile
}

class PetifiesHomeScreeen extends ConsumerStatefulWidget {
  const PetifiesHomeScreeen({super.key});

  @override
  ConsumerState<ConsumerStatefulWidget> createState() =>
      _PetifiesHomeScreeenState();
}

class _PetifiesHomeScreeenState extends ConsumerState<PetifiesHomeScreeen> {
  var _selectedTab = PetifiesHomeTabItem.explore;

  final _navigatorKeys = {
    PetifiesHomeTabItem.explore: GlobalKey<NavigatorState>(),
    PetifiesHomeTabItem.mypetifies: GlobalKey<NavigatorState>(),
    PetifiesHomeTabItem.newpetify: GlobalKey<NavigatorState>(),
    PetifiesHomeTabItem.notification: GlobalKey<NavigatorState>(),
    PetifiesHomeTabItem.profile: GlobalKey<NavigatorState>(),
  };

  @override
  void initState() {
    super.initState();
  }

  void _onItemTapped(int index) {
    setState(() {
      _selectedTab = PetifiesHomeTabItem.values[index];
    });
  }

  Widget _buildOffStageExploreNavigator() {
    return Offstage(
      offstage: _selectedTab != PetifiesHomeTabItem.explore,
      child: PetifiesExploreNavigator(
        navigatorKey: _navigatorKeys[PetifiesHomeTabItem.explore]!,
      ),
    );
  }

  Widget _buildOffStageMyPetifiesNavigator() {
    return Offstage(
      offstage: _selectedTab != PetifiesHomeTabItem.mypetifies,
      child: PetifiesDashboardNavigator(
        navigatorKey: _navigatorKeys[PetifiesHomeTabItem.mypetifies]!,
      ),
    );
  }

  Widget _buildOffStageNewPetifiesNavigator() {
    return Offstage(
      offstage: _selectedTab != PetifiesHomeTabItem.newpetify,
      child: Placeholder(),
    );
  }

  Widget _buildOffStageNotificationNavigator() {
    return Offstage(
      offstage: _selectedTab != PetifiesHomeTabItem.notification,
      child: Placeholder(),
    );
  }

  Widget _buildOffStageProfile() {
    return Offstage(
      offstage: _selectedTab != PetifiesHomeTabItem.profile,
      child: ProfileNavigator(
        navigatorKey: _navigatorKeys[PetifiesHomeTabItem.profile]!,
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
          _buildOffStageExploreNavigator(),
          _buildOffStageMyPetifiesNavigator(),
          _buildOffStageNewPetifiesNavigator(),
          _buildOffStageNotificationNavigator(),
          _buildOffStageProfile()
        ],
      ),
      bottomNavigationBar: PetifiesBottomNavBar(
        curPage: _selectedTab.index,
        onTapFunc: _onItemTapped,
      ),
    );
  }
}
