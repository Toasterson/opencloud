from django.db import models


class OS(models.Model):
    name = models.CharField(max_length=100)

    def __str__(self):
        return self.name


class RealmType(models.Model):
    """
    name: Name to be displayed
    template_path: path to template archive (ansible or similar format)
    Curent Types:
        Public: Everything is open in the Internet. Shared IP Stack with the Host. No Central Login
        Private: Apps are in an Seperate IP Range. The Host has a VPN Zone Launched by default connecting
                the Machines with eachother on other hosts or Clients from Corporate Networks. External Acces works via
                a Reverse Proxy.
        Corporate: The Realm has complete Identity Management SetUp with Private CA and
                    Reverse Proxy for external access
    """


class Realm(models.Model):
    """
    machines: [] Personal machines max one user per machine or dedicated servers/VM's
    users: [] Avatars that are able to log into this realms services
    admins: [] Avatars that are able to Administer this Realms Services
    type: What Type of realm to create (defines how to access services, firewall rules etc.)
    location: <Onsite/Europe/Africa...>
    instances: [] wordpress, krita, openerp, magento, desktops ...
    """
    LOCATIONS = (
        ('on', 'OnSite'),
        ('EU', 'Europe'),
        ('AU', 'Africa'),
    )
    name = models.CharField(max_length=100, blank=True, default='')
    created = models.DateTimeField(auto_now_add=True)
    creator = models.ForeignKey('auth.User', related_name='realms_owned', on_delete=models.PROTECT)
    type = models.ForeignKey('opcbase.RealmType', on_delete=models.PROTECT)
    users = models.ManyToManyField('opcbase.Avatar', related_name='member_of_realms')
    admins = models.ManyToManyField('opcbase.Avatar', related_name='admin_of_realms')
    REALMSTATES = (
        ('IN', 'Initializing'),
        ('WK', 'Working'),
        ('RE', 'Ready')
    )
    state = models.CharField(
        max_length=2,
        choices=REALMSTATES,
        default='IN'
    )
    location = models.CharField(
        max_length=2,
        choices=LOCATIONS,
        default='EU'
    )

    class Meta:
        ordering = ['name']
        permissions = (
            ('view_realm', 'View Realm'),
        )

    def __str__(self):
        return self.name


class Machine(models.Model):
    hostname = models.CharField(max_length=255, blank=True)
    created = models.DateTimeField(auto_now_add=True)
    ip = models.GenericIPAddressField(blank=True, null=True)
    domainname = models.CharField(max_length=255)
    os = models.ForeignKey('opcbase.OS', on_delete=models.PROTECT)
    realm = models.ForeignKey('opcbase.Realm', related_name='machines', on_delete=models.CASCADE)

    def __str__(self):
        return self.hostname + '.' + self.domainname

    class Meta:
        permissions = (
            ('view_machine', 'View Machine'),
        )


class Instance(models.Model):
    name = models.CharField(max_length=255)
    realm = models.ForeignKey('opcbase.Realm', related_name='instances', on_delete=models.CASCADE)
    app = models.ForeignKey('opcbase.App', on_delete=models.PROTECT)

    def __str__(self):
        return self.name

    class Meta:
        permissions = (
            ('view_instance', 'View Instance'),
        )


class App(models.Model):
    name = models.CharField(max_length=255)
    APPTYPES = (
        ('WEBA', 'WebApplication'),
        ('DESK', 'Desktop'),
        ('MOBI', 'Mobile')
    )
    type = models.CharField(
        max_length=4,
        choices=APPTYPES,
        default='WEBA'
    )
    maintainer = models.ForeignKey('auth.User', related_name='maintains_apps', on_delete=models.PROTECT)
    LANGUAGETYPES = (
        ('PHP', 'PHP'),
        ('PYTH', 'Python'),
        ('RUBY', 'Ruby on Rails'),
    )
    language = models.CharField(
        max_length=5,
        choices=LANGUAGETYPES
    )

    def __str__(self):
        return self.name


class Avatar(models.Model):
    user = models.ForeignKey('auth.User', related_name='avatar', on_delete=models.CASCADE, blank=True)
    firstname = models.CharField(max_length=255, blank=True)
    lastname = models.CharField(max_length=255, blank=True)
    loginname = models.CharField(max_length=255, blank=True)
    #TODO Profile Related Stuff
    #TODO Infrastructure Related fields
