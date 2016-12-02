from rest_framework import serializers
from .models import Realm, Machine, Instance, App, Avatar
from django.contrib.auth.models import User, Group


class UserSerializer(serializers.ModelSerializer):

    class Meta:
        model = User
        fields = '__all__'


class GroupSerializer(serializers.ModelSerializer):

    class Meta:
        model = Group
        fields = '__all__'


class RealmSerializer(serializers.ModelSerializer):
    creator = serializers.ReadOnlyField(source="creator.username")

    class Meta:
        model = Realm
        fields = '__all__'


class MachineSerializer(serializers.ModelSerializer):

    class Meta:
        model = Machine
        fields = '__all__'


class InstanceSerializer(serializers.ModelSerializer):

    class Meta:
        model = Instance
        fields = '__all__'


class AppSerializer(serializers.ModelSerializer):

    class Meta:
        model = App
        fields = '__all__'


class AvatarSerializer(serializers.ModelSerializer):

    class Meta:
        model = Avatar
        fields = '__all__'
