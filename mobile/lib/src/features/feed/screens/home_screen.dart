import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/providers/model_providers.dart';
import 'package:mobile/src/widgets/appbars/main_appbar.dart';
import 'package:mobile/src/widgets/bottom_nav_bars/main_bottom_nav_bar.dart';
import 'package:mobile/src/widgets/story_circle/story_circle.dart';

class HomeScreeen extends ConsumerWidget {
  const HomeScreeen({super.key});

  Widget get _stories {
    return SizedBox(
      height: 100,
      child: ListView(
        scrollDirection: Axis.horizontal,
        children: [
          StoryCircle(),
          StoryCircle(),
          StoryCircle(),
          StoryCircle(),
          StoryCircle(),
          StoryCircle(),
          StoryCircle(),
          StoryCircle(),
          StoryCircle(),
          StoryCircle(),
          StoryCircle(),
          StoryCircle(),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
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
      appBar: const MainAppBar(),
      body: Column(
        children: [
          Container(
            height: 5,
            color: Color.fromRGBO(10, 10, 10, 10),
          ),
          _stories,
          Container(
            height: 5,
            color: Color.fromRGBO(10, 10, 10, 10),
          ),
        ],
      ),
      bottomNavigationBar: MainButtomNavBar(),
    );
  }
}
