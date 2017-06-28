import mimetypes
import os

import sys
import uuid
from uuid import UUID

from django import conf
from django.core.files.uploadedfile import InMemoryUploadedFile
from rest_framework import filters
from rest_framework import mixins
from rest_framework.parsers import JSONParser, MultiPartParser
from opcbase import permissions
from opcbase.models import Realm, Machine, Instance, App
from opcbase.serializers import RealmSerializer, UserSerializer, GroupSerializer, MachineSerializer, \
    AppSerializer, InstanceSerializer
from django.contrib.auth.models import User, Group
from rest_framework import viewsets
from rest_framework.versioning import AcceptHeaderVersioning
from os import path
from opcbase.tasks import realm_add, realm_change, realm_destroy


class UserViewSet(viewsets.ReadOnlyModelViewSet):
    """
    This viewset automatically provides `list` and `detail` actions.
    """
    queryset = User.objects.all()
    serializer_class = UserSerializer
    versioning_class = AcceptHeaderVersioning


class GroupViewSet(viewsets.ModelViewSet):
    queryset = Group.objects.all()
    serializer_class = GroupSerializer
    versioning_class = AcceptHeaderVersioning


class RealmViewSet(viewsets.ModelViewSet):
    """
    This viewset automatically provides `list`, `create`, `retrieve`,
    `update` and `destroy` actions.
    """
    queryset = Realm.objects.all()
    serializer_class = RealmSerializer
    versioning_class = AcceptHeaderVersioning
    permission_classes = [permissions.RealmObjectPermissions]
    filter_backends = [filters.DjangoObjectPermissionsFilter]

    def perform_create(self, serializer):
        instance = serializer.save(creator=self.request.user)
        self.request.user.add_obj_perm('view_realm', instance)
        self.request.user.add_obj_perm('change_realm', instance)
        self.request.user.add_obj_perm('delete_realm', instance)
        realm_add({'id': instance.id})

    def perform_destroy(self, instance):
        realm_destroy({'id': instance.id})
        instance.delete()

    def perform_update(self, serializer):
        realm_change({'id': serializer.save().id})


class MachineViewSet(viewsets.ModelViewSet):
    queryset = Machine.objects.all()
    serializer_class = MachineSerializer
    versioning_class = AcceptHeaderVersioning
    filter_backends = [filters.DjangoObjectPermissionsFilter]

    def perform_create(self, serializer):
        instance = serializer.save()
        self.request.user.add_obj_perm('view_machine', instance)
        self.request.user.add_obj_perm('change_machine', instance)
        self.request.user.add_obj_perm('delete_machine', instance)


class InstanceViewSet(mixins.CreateModelMixin,
                      mixins.RetrieveModelMixin,
                      mixins.DestroyModelMixin,
                      mixins.ListModelMixin,
                      viewsets.GenericViewSet):
    queryset = Instance.objects.all()
    serializer_class = InstanceSerializer
    versioning_class = AcceptHeaderVersioning
    filter_backends = [filters.DjangoObjectPermissionsFilter]

    def perform_create(self, serializer):
        instance = serializer.save(name=uuid.uuid4())
        self.request.user.add_obj_perm('view_instance', instance)
        self.request.user.add_obj_perm('change_instance', instance)
        self.request.user.add_obj_perm('delete_instance', instance)


class AppViewSet(viewsets.ModelViewSet):
    queryset = App.objects.all()
    serializer_class = AppSerializer
    versioning_class = AcceptHeaderVersioning
    parser_classes = [JSONParser, MultiPartParser]

    def perform_destroy(self, instance):
        """
        Destroys the App archive can only be done by the maintainer
        TODO detroy method
        :param instance:
        :return:
        """
        pass

    def perform_create(self, serializer):
        if 'application/x-tar' in mimetypes.guess_type(serializer.validated_data['archive'].name):
            serializer.save(maintainer=self.request.user)

    def perform_update(self, serializer):
        if isinstance(serializer.validated_data['archive'], InMemoryUploadedFile):
            if 'application/x-tar' in mimetypes.guess_type(serializer.validated_data['archive'].name):
                serializer.save()
