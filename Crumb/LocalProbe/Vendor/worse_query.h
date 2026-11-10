//
//  worse_query.h
//  worse_query
//
//  Created by Taj C on 8/20/26.
//

/*
     worse_query - A sandbox escape PoC for iOS 26 and 27, again.
     Copyright (C) 2026 Taj C (forcequit)

     This program is free software: you can redistribute it and/or modify
     it under the terms of the GNU General Public License as published by
     the Free Software Foundation, either version 3 of the License, or
     (at your option) any later version.

     This program is distributed in the hope that it will be useful,
     but WITHOUT ANY WARRANTY; without even the implied warranty of
     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
     GNU General Public License for more details.

     You should have received a copy of the GNU General Public License
     along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

#ifndef worse_query_h
#define worse_query_h

#import <Foundation/Foundation.h>
#include <stdint.h>
#include <stdbool.h>

// Requires a caller-prepared Caches symlink; performs no file mutations.
int64_t worse_query_prepared(NSString *path, NSString *group, NSString *preparedCaches, bool alt_flags);
void worse_query_release(int64_t handle);

#endif /* worse_query_h */
