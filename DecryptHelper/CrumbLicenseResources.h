#import <Foundation/Foundation.h>

// Restores original license bytes from the signed pre-install snapshot. It
// never reads the installed SC_Info path or alters any Mach-O encryption flag.
BOOL crumb_copy_license_resource(NSURL *bundle, NSString *relative,
                                 NSDictionary *resource, NSURL *dest,
                                 NSError **error);

// Restores original non-Payload files/directories without overwriting the app.
BOOL crumb_copy_archive_resource(NSURL *bundle, NSString *relative,
                                 NSDictionary *resource, NSURL *dest,
                                 NSError **error);
BOOL crumb_create_archive_directory(NSURL *dest, NSString *relative,
                                    NSError **error);
