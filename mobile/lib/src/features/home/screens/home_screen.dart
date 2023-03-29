import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/features/home/navigators/home_navigator.dart';
import 'package:mobile/src/features/home/navigators/profile_navigator.dart';
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
      child: Placeholder(),
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
    // return Scaffold(
    //   body: Center(
    //     child: MediaView(
    //       imageUrls: [
    //         // "https://storage.googleapis.com/petifies-storage/2d9e369b-7a68-4a5a-8ed8-29d46866cc59-image_picker_36866E47-A1B5-4C13-A422-D8F0F2CEF9AE-70167-000013514860E757.jpg",
    //         "https://storage.googleapis.com/petifies-storage/05560c77-c496-409d-84ef-0b595b29b58b-image_picker_53C8CFB3-8C41-415D-BB95-5D4C43F6DBD1-91911-000018AE78699AE1.jpg",
    //         // "https://storage.googleapis.com/petifies-storage/05560c77-c496-409d-84ef-0b595b29b58b-image_picker_53C8CFB3-8C41-415D-BB95-5D4C43F6DBD1-91911-000018AE78699AE1.jpg",
    //         // "https://storage.googleapis.com/petifies-storage/05560c77-c496-409d-84ef-0b595b29b58b-image_picker_53C8CFB3-8C41-415D-BB95-5D4C43F6DBD1-91911-000018AE78699AE1.jpg"
    //       ],
    //       videoUrls: [
    //         "https://dms.licdn.com/playlist/C5605AQFpDDHvqJofDA/mp4-720p-30fp-crf28/0/1679796920321?e=1680462000&v=beta&t=r0IU_2keGwXxhfNCX0cXpaapy61Df1WAYKQg6L165Co",
    //         "https://storage.googleapis.com/gtv-videos-bucket/sample/BigBuckBunny.mp4",
    //         "https://storage.googleapis.com/gtv-videos-bucket/sample/BigBuckBunny.mp4",
    //         // "https://storage.googleapis.com/gtv-videos-bucket/sample/BigBuckBunny.mp4"
    //       ],
    //     ),
    //   ),
    // );
    // return MediaFullPageScreen(
    //     mediaUrl:
    //         "https://storage.googleapis.com/petifies-storage/2d9e369b-7a68-4a5a-8ed8-29d46866cc59-image_picker_36866E47-A1B5-4C13-A422-D8F0F2CEF9AE-70167-000013514860E757.jpg",
    //     isMediaImage: true);

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
