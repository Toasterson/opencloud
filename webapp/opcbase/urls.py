"""webapp URL Configuration

The `urlpatterns` list routes URLs to views. For more information please see:
    https://docs.djangoproject.com/en/1.10/topics/http/urls/
Examples:
Function views
    1. Add an import:  from my_app import views
    2. Add a URL to urlpatterns:  url(r'^$', views.home, name='home')
Class-based views
    1. Add an import:  from other_app.views import Home
    2. Add a URL to urlpatterns:  url(r'^$', Home.as_view(), name='home')
Including another URLconf
    1. Import the include() function: from django.conf.urls import url, include
    2. Add a URL to urlpatterns:  url(r'^blog/', include('blog.urls'))
"""
from django.conf.urls import url, include
from rest_framework.schemas import get_schema_view
from opcbase.views import RealmViewSet, UserViewSet, MachineViewSet, InstanceViewSet, AppViewSet, GroupViewSet
from rest_framework.routers import DefaultRouter

# Create a router and register our viewsets with it.
schema_view = get_schema_view(title='OpenCloud API')
router = DefaultRouter()
router.register(r'realms', RealmViewSet)
router.register(r'users', UserViewSet)
router.register(r'groups', GroupViewSet)
router.register(r'machines', MachineViewSet)
router.register(r'instances', InstanceViewSet)
router.register(r'apps', AppViewSet)


urlpatterns = [
    url('^schema/$', schema_view),
    url(r'^', include(router.urls)),
]
